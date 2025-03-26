package fi.ehin.service;

import fi.ehin.external.NordPoolClient;
import io.quarkus.logging.Log;
import jakarta.inject.Singleton;
import java.time.LocalDate;
import org.eclipse.microprofile.rest.client.inject.RestClient;

@Singleton
public class PricesService {

  @RestClient
  NordPoolClient nordPoolClient;

  public NordPoolClient.PriceDataResponse GetTomorrowsPrices() {
    final var prices = nordPoolClient.getPricesForDate(
      LocalDate.now().plusDays(1),
      "DayAhead",
      "FI",
      "EUR"
    );
    if (invalidPrices(prices)) {
      return null;
    }
    return prices;
  }

  public NordPoolClient.PriceDataResponse GetPrices(final LocalDate date) {
    final var prices = nordPoolClient.getPricesForDate(
      date,
      "DayAhead",
      "FI",
      "EUR"
    );
    if (invalidPrices(prices)) {
      return null;
    }
    return prices;
  }

  private boolean invalidPrices(final NordPoolClient.PriceDataResponse prices) {
    if (prices == null) {
      Log.errorf("Expected to find prices but was null");
      return true;
    }
    if (!"DayAhead".equals(prices.market())) {
      Log.errorf("Expected market DayAhead but got %s", prices.market());
      return true;
    }
    if (!"EUR".equals(prices.currency())) {
      Log.errorf("Expected currency EUR but got %s", prices.market());
      return true;
    }
    if (prices.areaStates() == null || prices.areaStates().isEmpty()) {
      Log.errorf("Expected areaStates to not be empty");
      return true;
    }
    final var fiState = prices
      .areaStates()
      .stream()
      .filter(s -> s.areas() != null && s.areas().contains("FI"))
      .findFirst();
    if (fiState.isEmpty()) {
      Log.errorf("Couldn't find FI area from area states");
      return true;
    }
    if (!"Final".equals(fiState.get().state())) {
      Log.errorf("Expected state Final but got %s", fiState.get().state());
      return true;
    }
    return false;
  }
}
