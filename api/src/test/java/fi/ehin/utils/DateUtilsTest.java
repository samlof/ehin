package fi.ehin.utils;

import io.quarkus.test.junit.QuarkusTest;
import java.time.OffsetDateTime;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

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
}
