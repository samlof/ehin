package fi.ehin.external;

import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.QueryParam;
import java.math.BigDecimal;
import java.time.LocalDate;
import java.time.OffsetDateTime;
import java.util.List;
import org.eclipse.microprofile.rest.client.inject.RegisterRestClient;

// https://dataportal-api.nordpoolgroup.com/api/DayAheadPrices?date=2025-09-30&market=DayAhead&deliveryArea=FI&currency=EUR

@RegisterRestClient(
  baseUri = "https://dataportal-api.nordpoolgroup.com",
  configKey = "nordpool-api"
)
public interface NordPoolClient {
  record PriceDataResponse(
    LocalDate deliveryDateCET,
    int version,
    List<String> deliveryAreas,
    OffsetDateTime updatedAt,
    String market,
    String currency,
    List<PriceDataResponseEntry> multiAreaEntries,
    List<PriceDataResponseState> areaStates,
    List<PriceDataResponseAverage> areaAverages
  ) {}

  record PriceDataResponseEntry(
    OffsetDateTime deliveryStart,
    OffsetDateTime deliveryEnd,
    PriceDataResponseData entryPerArea
  ) {}

  record PriceDataResponseData(BigDecimal FI) {}

  record PriceDataResponseAverage(String areaCode, BigDecimal price) {}

  record PriceDataResponseState(String state, List<String> areas) {}

  @Path("/api/DayAheadPrices")
  @GET
  PriceDataResponse getPricesForDate(
    @QueryParam("date") LocalDate date,
    @QueryParam("market") String market,
    @QueryParam("deliveryArea") String deliveryArea,
    @QueryParam("currency") String currency
  );
}
