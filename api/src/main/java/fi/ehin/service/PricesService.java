package fi.ehin.service;

import fi.ehin.external.NordPoolClient;
import jakarta.inject.Singleton;
import java.time.LocalDate;
import org.eclipse.microprofile.rest.client.inject.RestClient;

@Singleton
public class PricesService {

  @RestClient
  NordPoolClient nordPoolClient;

  public NordPoolClient.PriceDataResponse GetTomorrowsPrices() {
    return nordPoolClient.getPricesForDate(
      LocalDate.now().plusDays(1),
      "DayAhead",
      "FI",
      "EUR"
    );
  }

  public NordPoolClient.PriceDataResponse GetPrices(final LocalDate date) {
    return nordPoolClient.getPricesForDate(date, "DayAhead", "FI", "EUR");
  }
}
