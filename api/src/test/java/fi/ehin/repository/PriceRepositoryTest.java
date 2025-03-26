package fi.ehin.repository;

import static org.hamcrest.CoreMatchers.hasItem;
import static org.hamcrest.MatcherAssert.assertThat;

import fi.ehin.external.NordPoolClient;
import io.quarkus.test.junit.QuarkusTest;
import jakarta.inject.Inject;
import java.math.BigDecimal;
import java.time.OffsetDateTime;
import java.util.List;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

@QuarkusTest
public class PriceRepositoryTest {

  @Inject
  PriceRepository priceRepository;

  @Test
  void testSelect1Works() {
    final var one = priceRepository.select1();
    Assertions.assertEquals(1, one);
  }

  @Test
  void testInsertAndGetPrices() {
    priceRepository.insertPrices(
      List.of(
        new NordPoolClient.PriceDataResponseEntry(
          OffsetDateTime.now().withHour(0),
          OffsetDateTime.now().withHour(1),
          new NordPoolClient.PriceDataResponseData(BigDecimal.valueOf(1.00))
        )
      )
    );

    final var rows = priceRepository.getPrices(
      OffsetDateTime.now().withHour(0).minusMinutes(1),
      OffsetDateTime.now().withHour(1).plusMinutes(1)
    );

    Assertions.assertFalse(rows.isEmpty(), "expected to find rows");
  }
}
