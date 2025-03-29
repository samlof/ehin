package fi.ehin.utils;

import io.quarkus.test.junit.QuarkusTest;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.OffsetDateTime;
import java.time.OffsetTime;
import java.time.ZoneOffset;

import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

import static fi.ehin.utils.DateUtils.HELSINKI_ZONE;


@QuarkusTest
public class DateUtilsTest {

  @Test
  void testTimeUntilUtc() {
    final var seconds = DateUtils.secondsUntil12(
      OffsetDateTime.now().withHour(11).withMinute(0).withSecond(0).withNano(0)
    );
    Assertions.assertEquals(3600, seconds);

    final var seconds2 = DateUtils.secondsUntil12(
      OffsetDateTime.now().withHour(11).withMinute(1).withSecond(0).withNano(0)
    );
    Assertions.assertEquals(3540, seconds2);

    final var seconds3 = DateUtils.secondsUntil12(
      OffsetDateTime.now().withHour(11).withMinute(1).withSecond(5).withNano(0)
    );
    Assertions.assertEquals(3535, seconds3);
  }

  @Test
  void testGmtString() {
    Assertions.assertEquals("Sat, 29 Mar 2025 11:57:00 GMT", DateUtils.getGmtStringForCache(LocalDate.of(2025, 3, 29), 0));
    Assertions.assertEquals("Sat, 22 Feb 2025 11:57:25 GMT", DateUtils.getGmtStringForCache(LocalDate.of(2025, 2, 22), 25));
    Assertions.assertEquals("Thu, 29 May 2025 11:57:50 GMT", DateUtils.getGmtStringForCache(LocalDate.of(2025, 5, 29), 50));
  }


}
