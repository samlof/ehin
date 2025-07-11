package fi.ehin.resource;

import static fi.ehin.utils.RequestUtils.*;
import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.*;

import fi.ehin.repository.PriceRepository;
import fi.ehin.service.DateService;
import fi.ehin.utils.DateUtils;
import io.quarkus.test.InjectMock;
import io.quarkus.test.common.http.TestHTTPEndpoint;
import io.quarkus.test.junit.QuarkusTest;
import java.math.BigDecimal;
import java.time.Instant;
import java.time.OffsetDateTime;
import java.time.ZonedDateTime;
import java.util.List;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

@QuarkusTest
@TestHTTPEndpoint(PriceResource.class)
class PriceResourceTest {

  @InjectMock
  PriceRepository priceRepository;

  @InjectMock
  DateService dateService;

  @Test
  void testGetTomorrowsEntries() {
    Mockito.when(
      priceRepository.getPrices(Mockito.any(), Mockito.any())
    ).thenReturn(makeEntriesWithTomorrow());
    Mockito.when(dateService.now()).thenReturn(
      ZonedDateTime.of(
        2025,
        3,
        29,
        13,
        1,
        1,
        1,
        DateUtils.HELSINKI_ZONE
      ).toOffsetDateTime()
    );

    given()
      .when()
      .get("/prices/2025-03-29")
      .then()
      .statusCode(200)
      .header(CACHE_CONTROL_HEADER, CACHE_LONG)
      .body(containsString("2025-03-30T21:00:00Z"));
  }

  @Test
  void testGetTomorrowsEntriesMissingAtOne() {
    Mockito.when(
      priceRepository.getPrices(Mockito.any(), Mockito.any())
    ).thenReturn(makeEntriesWithoutAll());
    Mockito.when(dateService.now()).thenReturn(
      ZonedDateTime.of(
        2025,
        3,
        29,
        1,
        1,
        1,
        1,
        DateUtils.HELSINKI_ZONE
      ).toOffsetDateTime()
    );

    given()
      .when()
      .get("/prices/2025-03-29")
      .then()
      .statusCode(200)
      .header(
        EXPIRES_HEADER,
        allOf(startsWith("Sat, 29 Mar 2025 11:57:"), endsWith("GMT"))
      )
      .header(CACHE_CONTROL_HEADER, CACHE_VAR)
      .body(not(containsString("2025-03-30T21:00:00Z")));
  }

  @Test
  void testGetTomorrowsEntriesMissing() {
    Mockito.when(
      priceRepository.getPrices(Mockito.any(), Mockito.any())
    ).thenReturn(makeEntriesWithoutAll());
    Mockito.when(dateService.now()).thenReturn(
      ZonedDateTime.of(
        2025,
        3,
        29,
        10,
        1,
        1,
        1,
        DateUtils.HELSINKI_ZONE
      ).toOffsetDateTime()
    );

    given()
      .when()
      .get("/prices/2025-03-29")
      .then()
      .statusCode(200)
      .header(
        EXPIRES_HEADER,
        allOf(startsWith("Sat, 29 Mar 2025 11:57:"), endsWith("GMT"))
      )
      .header(CACHE_CONTROL_HEADER, CACHE_VAR)
      .body(not(containsString("2025-03-30T21:00:00Z")));
  }

  @Test
  void testGetTomorrowsEntriesMissingEarlier() {
    Mockito.when(
      priceRepository.getPrices(Mockito.any(), Mockito.any())
    ).thenReturn(makeEntriesWithoutAll());
    Mockito.when(dateService.now()).thenReturn(
      ZonedDateTime.of(
        2025,
        3,
        29,
        13,
        58,
        1,
        1,
        DateUtils.HELSINKI_ZONE
      ).toOffsetDateTime()
    );

    given()
      .when()
      .get("/prices/2025-03-29")
      .then()
      .statusCode(200)
      .header(CACHE_CONTROL_HEADER, CACHE_VAR + ", max-age=60")
      .header(EXPIRES_HEADER, nullValue())
      .body(not(containsString("2025-03-30T21:00:00Z")));
  }

  private static PriceRepository.PriceHistoryEntry makeEntry(
    final String end,
    final String start,
    final double price
  ) {
    return new PriceRepository.PriceHistoryEntry(
      BigDecimal.valueOf(price),
      OffsetDateTime.parse(start),
      OffsetDateTime.parse(end)
    );
  }

  private static List<
    PriceRepository.PriceHistoryEntry
  > makeEntriesWithoutAll() {
    return List.of(
      makeEntry("2025-03-27T23:00:00Z", "2025-03-27T22:00:00Z", 3.4),
      makeEntry("2025-03-28T00:00:00Z", "2025-03-27T23:00:00Z", 2.99),
      makeEntry("2025-03-28T01:00:00Z", "2025-03-28T00:00:00Z", 2.23),
      makeEntry("2025-03-28T02:00:00Z", "2025-03-28T01:00:00Z", 1.29),
      makeEntry("2025-03-28T03:00:00Z", "2025-03-28T02:00:00Z", 1.28),
      makeEntry("2025-03-28T04:00:00Z", "2025-03-28T03:00:00Z", 4.87),
      makeEntry("2025-03-28T05:00:00Z", "2025-03-28T04:00:00Z", 9.31),
      makeEntry("2025-03-28T06:00:00Z", "2025-03-28T05:00:00Z", 11.17),
      makeEntry("2025-03-28T07:00:00Z", "2025-03-28T06:00:00Z", 12.14),
      makeEntry("2025-03-28T08:00:00Z", "2025-03-28T07:00:00Z", 12.65),
      makeEntry("2025-03-28T09:00:00Z", "2025-03-28T08:00:00Z", 11.41),
      makeEntry("2025-03-28T10:00:00Z", "2025-03-28T09:00:00Z", 6.89),
      makeEntry("2025-03-28T11:00:00Z", "2025-03-28T10:00:00Z", 1.3),
      makeEntry("2025-03-28T12:00:00Z", "2025-03-28T11:00:00Z", -0.01),
      makeEntry("2025-03-28T13:00:00Z", "2025-03-28T12:00:00Z", -0.02),
      makeEntry("2025-03-28T14:00:00Z", "2025-03-28T13:00:00Z", -0.01),
      makeEntry("2025-03-28T15:00:00Z", "2025-03-28T14:00:00Z", 0.01),
      makeEntry("2025-03-28T16:00:00Z", "2025-03-28T15:00:00Z", 0.66),
      makeEntry("2025-03-28T17:00:00Z", "2025-03-28T16:00:00Z", 1.23),
      makeEntry("2025-03-28T18:00:00Z", "2025-03-28T17:00:00Z", 1.25),
      makeEntry("2025-03-28T19:00:00Z", "2025-03-28T18:00:00Z", 0.41),
      makeEntry("2025-03-28T20:00:00Z", "2025-03-28T19:00:00Z", 0.0),
      makeEntry("2025-03-28T21:00:00Z", "2025-03-28T20:00:00Z", 0.01),
      makeEntry("2025-03-28T22:00:00Z", "2025-03-28T21:00:00Z", 0.0),
      makeEntry("2025-03-28T23:00:00Z", "2025-03-28T22:00:00Z", -0.01),
      makeEntry("2025-03-29T00:00:00Z", "2025-03-28T23:00:00Z", 0.0),
      makeEntry("2025-03-29T01:00:00Z", "2025-03-29T00:00:00Z", 0.0),
      makeEntry("2025-03-29T02:00:00Z", "2025-03-29T01:00:00Z", 0.0),
      makeEntry("2025-03-29T03:00:00Z", "2025-03-29T02:00:00Z", 0.0),
      makeEntry("2025-03-29T04:00:00Z", "2025-03-29T03:00:00Z", -0.01),
      makeEntry("2025-03-29T05:00:00Z", "2025-03-29T04:00:00Z", 0.01),
      makeEntry("2025-03-29T06:00:00Z", "2025-03-29T05:00:00Z", 0.91),
      makeEntry("2025-03-29T07:00:00Z", "2025-03-29T06:00:00Z", 1.26),
      makeEntry("2025-03-29T08:00:00Z", "2025-03-29T07:00:00Z", 1.91),
      makeEntry("2025-03-29T09:00:00Z", "2025-03-29T08:00:00Z", 1.24),
      makeEntry("2025-03-29T10:00:00Z", "2025-03-29T09:00:00Z", 0.01),
      makeEntry("2025-03-29T11:00:00Z", "2025-03-29T10:00:00Z", 0.0),
      makeEntry("2025-03-29T12:00:00Z", "2025-03-29T11:00:00Z", 0.0),
      makeEntry("2025-03-29T13:00:00Z", "2025-03-29T12:00:00Z", -0.02),
      makeEntry("2025-03-29T14:00:00Z", "2025-03-29T13:00:00Z", 0.73),
      makeEntry("2025-03-29T15:00:00Z", "2025-03-29T14:00:00Z", 2.59),
      makeEntry("2025-03-29T16:00:00Z", "2025-03-29T15:00:00Z", 3.03),
      makeEntry("2025-03-29T17:00:00Z", "2025-03-29T16:00:00Z", 3.39),
      makeEntry("2025-03-29T18:00:00Z", "2025-03-29T17:00:00Z", 3.9),
      makeEntry("2025-03-29T19:00:00Z", "2025-03-29T18:00:00Z", 3.82),
      makeEntry("2025-03-29T20:00:00Z", "2025-03-29T19:00:00Z", 3.84),
      makeEntry("2025-03-29T21:00:00Z", "2025-03-29T20:00:00Z", 4.0),
      makeEntry("2025-03-29T22:00:00Z", "2025-03-29T21:00:00Z", 4.2),
      makeEntry("2025-03-29T23:00:00Z", "2025-03-29T22:00:00Z", 4.79)
    );
  }

  private static List<
    PriceRepository.PriceHistoryEntry
  > makeEntriesWithTomorrow() {
    return List.of(
      makeEntry("2025-03-27T23:00:00Z", "2025-03-27T22:00:00Z", 3.4),
      makeEntry("2025-03-28T00:00:00Z", "2025-03-27T23:00:00Z", 2.99),
      makeEntry("2025-03-28T01:00:00Z", "2025-03-28T00:00:00Z", 2.23),
      makeEntry("2025-03-28T02:00:00Z", "2025-03-28T01:00:00Z", 1.29),
      makeEntry("2025-03-28T03:00:00Z", "2025-03-28T02:00:00Z", 1.28),
      makeEntry("2025-03-28T04:00:00Z", "2025-03-28T03:00:00Z", 4.87),
      makeEntry("2025-03-28T05:00:00Z", "2025-03-28T04:00:00Z", 9.31),
      makeEntry("2025-03-28T06:00:00Z", "2025-03-28T05:00:00Z", 11.17),
      makeEntry("2025-03-28T07:00:00Z", "2025-03-28T06:00:00Z", 12.14),
      makeEntry("2025-03-28T08:00:00Z", "2025-03-28T07:00:00Z", 12.65),
      makeEntry("2025-03-28T09:00:00Z", "2025-03-28T08:00:00Z", 11.41),
      makeEntry("2025-03-28T10:00:00Z", "2025-03-28T09:00:00Z", 6.89),
      makeEntry("2025-03-28T11:00:00Z", "2025-03-28T10:00:00Z", 1.3),
      makeEntry("2025-03-28T12:00:00Z", "2025-03-28T11:00:00Z", -0.01),
      makeEntry("2025-03-28T13:00:00Z", "2025-03-28T12:00:00Z", -0.02),
      makeEntry("2025-03-28T14:00:00Z", "2025-03-28T13:00:00Z", -0.01),
      makeEntry("2025-03-28T15:00:00Z", "2025-03-28T14:00:00Z", 0.01),
      makeEntry("2025-03-28T16:00:00Z", "2025-03-28T15:00:00Z", 0.66),
      makeEntry("2025-03-28T17:00:00Z", "2025-03-28T16:00:00Z", 1.23),
      makeEntry("2025-03-28T18:00:00Z", "2025-03-28T17:00:00Z", 1.25),
      makeEntry("2025-03-28T19:00:00Z", "2025-03-28T18:00:00Z", 0.41),
      makeEntry("2025-03-28T20:00:00Z", "2025-03-28T19:00:00Z", 0.0),
      makeEntry("2025-03-28T21:00:00Z", "2025-03-28T20:00:00Z", 0.01),
      makeEntry("2025-03-28T22:00:00Z", "2025-03-28T21:00:00Z", 0.0),
      makeEntry("2025-03-28T23:00:00Z", "2025-03-28T22:00:00Z", -0.01),
      makeEntry("2025-03-29T00:00:00Z", "2025-03-28T23:00:00Z", 0.0),
      makeEntry("2025-03-29T01:00:00Z", "2025-03-29T00:00:00Z", 0.0),
      makeEntry("2025-03-29T02:00:00Z", "2025-03-29T01:00:00Z", 0.0),
      makeEntry("2025-03-29T03:00:00Z", "2025-03-29T02:00:00Z", 0.0),
      makeEntry("2025-03-29T04:00:00Z", "2025-03-29T03:00:00Z", -0.01),
      makeEntry("2025-03-29T05:00:00Z", "2025-03-29T04:00:00Z", 0.01),
      makeEntry("2025-03-29T06:00:00Z", "2025-03-29T05:00:00Z", 0.91),
      makeEntry("2025-03-29T07:00:00Z", "2025-03-29T06:00:00Z", 1.26),
      makeEntry("2025-03-29T08:00:00Z", "2025-03-29T07:00:00Z", 1.91),
      makeEntry("2025-03-29T09:00:00Z", "2025-03-29T08:00:00Z", 1.24),
      makeEntry("2025-03-29T10:00:00Z", "2025-03-29T09:00:00Z", 0.01),
      makeEntry("2025-03-29T11:00:00Z", "2025-03-29T10:00:00Z", 0.0),
      makeEntry("2025-03-29T12:00:00Z", "2025-03-29T11:00:00Z", 0.0),
      makeEntry("2025-03-29T13:00:00Z", "2025-03-29T12:00:00Z", -0.02),
      makeEntry("2025-03-29T14:00:00Z", "2025-03-29T13:00:00Z", 0.73),
      makeEntry("2025-03-29T15:00:00Z", "2025-03-29T14:00:00Z", 2.59),
      makeEntry("2025-03-29T16:00:00Z", "2025-03-29T15:00:00Z", 3.03),
      makeEntry("2025-03-29T17:00:00Z", "2025-03-29T16:00:00Z", 3.39),
      makeEntry("2025-03-29T18:00:00Z", "2025-03-29T17:00:00Z", 3.9),
      makeEntry("2025-03-29T19:00:00Z", "2025-03-29T18:00:00Z", 3.82),
      makeEntry("2025-03-29T20:00:00Z", "2025-03-29T19:00:00Z", 3.84),
      makeEntry("2025-03-29T21:00:00Z", "2025-03-29T20:00:00Z", 4.0),
      makeEntry("2025-03-29T22:00:00Z", "2025-03-29T21:00:00Z", 4.2),
      makeEntry("2025-03-29T23:00:00Z", "2025-03-29T22:00:00Z", 4.79),
      makeEntry("2025-03-30T00:00:00Z", "2025-03-29T23:00:00Z", 5.96),
      makeEntry("2025-03-30T01:00:00Z", "2025-03-30T00:00:00Z", 6.92),
      makeEntry("2025-03-30T02:00:00Z", "2025-03-30T01:00:00Z", 4.79),
      makeEntry("2025-03-30T03:00:00Z", "2025-03-30T02:00:00Z", 3.9),
      makeEntry("2025-03-30T04:00:00Z", "2025-03-30T03:00:00Z", 3.9),
      makeEntry("2025-03-30T05:00:00Z", "2025-03-30T04:00:00Z", 3.84),
      makeEntry("2025-03-30T06:00:00Z", "2025-03-30T05:00:00Z", 3.93),
      makeEntry("2025-03-30T07:00:00Z", "2025-03-30T06:00:00Z", 10.17),
      makeEntry("2025-03-30T08:00:00Z", "2025-03-30T07:00:00Z", 35.0),
      makeEntry("2025-03-30T09:00:00Z", "2025-03-30T08:00:00Z", 61.37),
      makeEntry("2025-03-30T10:00:00Z", "2025-03-30T09:00:00Z", 59.69),
      makeEntry("2025-03-30T11:00:00Z", "2025-03-30T10:00:00Z", 43.3),
      makeEntry("2025-03-30T12:00:00Z", "2025-03-30T11:00:00Z", 55.16),
      makeEntry("2025-03-30T13:00:00Z", "2025-03-30T12:00:00Z", 99.96),
      makeEntry("2025-03-30T14:00:00Z", "2025-03-30T13:00:00Z", 87.03),
      makeEntry("2025-03-30T15:00:00Z", "2025-03-30T14:00:00Z", 106.75),
      makeEntry("2025-03-30T16:00:00Z", "2025-03-30T15:00:00Z", 129.77),
      makeEntry("2025-03-30T17:00:00Z", "2025-03-30T16:00:00Z", 167.96),
      makeEntry("2025-03-30T18:00:00Z", "2025-03-30T17:00:00Z", 158.12),
      makeEntry("2025-03-30T19:00:00Z", "2025-03-30T18:00:00Z", 86.42),
      makeEntry("2025-03-30T20:00:00Z", "2025-03-30T19:00:00Z", 91.66),
      makeEntry("2025-03-30T21:00:00Z", "2025-03-30T20:00:00Z", 71.04),
      makeEntry("2025-03-30T22:00:00Z", "2025-03-30T21:00:00Z", 34.99)
    );
  }
}
