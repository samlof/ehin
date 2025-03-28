package fi.ehin.utils;

import io.quarkus.test.junit.QuarkusTest;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

@QuarkusTest
public class DateUtilsTest {

  @Test
  void testTimeUntilUtc() {
    final var seconds = DateUtils.secondsUntil12(
      OffsetDateTime.now().withHour(11).withMinute(0).withSecond(0)
    );
    Assertions.assertEquals(3600, seconds);
  }
}
