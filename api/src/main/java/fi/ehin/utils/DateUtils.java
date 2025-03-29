package fi.ehin.utils;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.OffsetDateTime;
import java.time.OffsetTime;
import java.time.ZoneId;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.time.temporal.ChronoUnit;
import java.util.random.RandomGenerator;

public final class DateUtils {

  private DateUtils() {}

  public static final ZoneId HELSINKI_ZONE = ZoneId.of("Europe/Helsinki");

  public static long secondsUntil12(final OffsetDateTime fromTime) {
    final var fromTimeWithUtc = fromTime.withOffsetSameInstant(ZoneOffset.UTC);
    return fromTimeWithUtc.until(
      fromTimeWithUtc.withHour(12).withMinute(0).withSecond(0).withNano(0),
      ChronoUnit.SECONDS
    );
  }

  public static long secondsUntil12Utc(final OffsetDateTime now) {
    return secondsUntil12(now.withNano(0));
  }

  public static String getGmtStringForCache(final OffsetDateTime now) {
    final var generator = RandomGenerator.getDefault();

    final var seconds = generator.nextInt(0, 59);
    return getGmtStringForCache(now.toLocalDate(), seconds);
  }

  public static String getGmtStringForCache(
    final LocalDate date,
    final int seconds
  ) {
    var dateInUtc = date.atTime(
      OffsetTime.of(
        0,
        0,
        0,
        0,
        ZoneOffset.UTC.getRules().getOffset(LocalDateTime.now())
      )
    );
    dateInUtc = dateInUtc
      .withHour(11)
      .withMinute(57)
      .withSecond(seconds)
      .withNano(0);
    return dateInUtc.format(DateTimeFormatter.RFC_1123_DATE_TIME);
  }
}
