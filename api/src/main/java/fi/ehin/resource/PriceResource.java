package fi.ehin.resource;

import fi.ehin.repository.PriceRepository;
import fi.ehin.service.PricesService;
import fi.ehin.utils.DateUtils;
import io.quarkus.logging.Log;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.PathParam;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.QueryParam;
import jakarta.ws.rs.core.MediaType;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.OffsetTime;
import java.time.ZoneId;
import java.util.List;
import org.eclipse.microprofile.config.inject.ConfigProperty;
import org.jboss.resteasy.reactive.RestResponse;

@Path("/api")
public class PriceResource {

  private final PriceRepository priceRepository;
  private final PricesService pricesService;

  private static final String CACHE_VAR = "public, max-age=";
  private static final String CACHE_LONG =
    "public, max-age=" + 60 * 60 * 168 + ", immutable";

  private static final ZoneId HELSINKI_ZONE = ZoneId.of("Europe/Helsinki");

  @ConfigProperty(name = "update-prices.password")
  String updatePricesPassword;

  public PriceResource(
    PriceRepository priceRepository,
    PricesService pricesService
  ) {
    this.priceRepository = priceRepository;
    this.pricesService = pricesService;
  }

  @GET
  @Produces(MediaType.APPLICATION_JSON)
  @Path("/prices/{date}")
  public RestResponse<List<PriceRepository.PriceHistoryEntry>> getPastPrices(
    @PathParam("date") final LocalDate date
  ) {
    Log.infof("Fetching prices for %s", date);

    final var dateWithTime = date.atTime(
      OffsetTime.of(
        0,
        0,
        0,
        0,
        HELSINKI_ZONE.getRules().getOffset(LocalDateTime.now())
      )
    );

    // The latest price is the day after's 0-1 prices. So plusDays 3 and for checking plusDays 2
    final var prices = priceRepository.getPrices(
      dateWithTime.minusDays(1),
      dateWithTime.plusDays(3)
    );
    String cacheString = CACHE_LONG; // Long cache if already have all prices
    final var newestPriceDate = prices.getLast().deliveryStart().toLocalDate();
    if (!newestPriceDate.equals(date.plusDays(2))) {
      // Wait seconds until 11:57:30 but at least 60 seconds
      final var waitTime = Math.max(DateUtils.secondsUntil12Utc() - 60*2+30, 60);
      Log.infof("Returning %d seconds cache for prices on %s", waitTime, date);
      cacheString = CACHE_VAR + waitTime;
    }
    return RestResponse.ResponseBuilder.ok(prices)
      .header("Cache-Control", cacheString)
      .build();
  }

  public record UpdatePricesResponse(boolean done) {}

  @GET
  @Produces(MediaType.APPLICATION_JSON)
  @Path("/update-prices")
  public RestResponse<UpdatePricesResponse> updatePrices(
    @QueryParam("p") final String password
  ) {
    if (!updatePricesPassword.equals(password)) {
      return RestResponse.notFound();
    }
    Log.infof("Updating prices for tomorrow %s", LocalDate.now().plusDays(1));
    final var prices = pricesService.GetTomorrowsPrices();
    if (prices == null) {
      return RestResponse.ok(new UpdatePricesResponse(false));
    }
    priceRepository.insertPrices(prices.multiAreaEntries());

    return RestResponse.ok(new UpdatePricesResponse(true));
  }

  @GET
  @Produces(MediaType.APPLICATION_JSON)
  @Path("/update-prices/{date}")
  public RestResponse<UpdatePricesResponse> updatePricesForDate(
    @PathParam("date") final LocalDate date,
    @QueryParam("p") final String password
  ) {
    if (!updatePricesPassword.equals(password)) {
      return RestResponse.notFound();
    }
    Log.infof("Updating prices for %s", date);
    final var prices = pricesService.GetPrices(date);
    priceRepository.insertPrices(prices.multiAreaEntries());

    return RestResponse.ok(new UpdatePricesResponse(true));
  }
}
