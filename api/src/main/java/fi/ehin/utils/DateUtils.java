package fi.ehin.utils;

import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.time.temporal.ChronoUnit;

public final class DateUtils {

  private DateUtils() {}

  public static long secondsUntil12(final OffsetDateTime fromTime) {
    return fromTime.until(
      fromTime.withHour(12).withMinute(0).withSecond(0).withNano(0),
      ChronoUnit.SECONDS
    );
  }

  public static long secondsUntil12Utc() {
    return secondsUntil12(OffsetDateTime.now(ZoneOffset.UTC));
  }
}
