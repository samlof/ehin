package fi.ehin.resource;

import fi.ehin.repository.PriceRepository;
import fi.ehin.service.DateService;
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
import java.time.LocalTime;
import java.time.OffsetDateTime;
import java.time.OffsetTime;
import java.time.ZoneOffset;
import java.util.List;
import org.eclipse.microprofile.config.inject.ConfigProperty;
import org.jboss.resteasy.reactive.RestResponse;

import static fi.ehin.utils.DateUtils.HELSINKI_ZONE;
import static fi.ehin.utils.RequestUtils.*;

@Path("/api")
public class PriceResource {

  private final PriceRepository priceRepository;
  private final PricesService pricesService;
  private final DateService dateService;

  @ConfigProperty(name = "update-prices.password")
  String updatePricesPassword;

  public PriceResource(
          PriceRepository priceRepository,
          PricesService pricesService, DateService dateService
  ) {
    this.priceRepository = priceRepository;
    this.pricesService = pricesService;
      this.dateService = dateService;
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
    String expiresValue = null;
    final var newestPriceDate = prices.getLast().deliveryStart().toLocalDate();
    if (!newestPriceDate.isAfter( date)) {
      final var pricesUpdateTime = OffsetDateTime.of(date, LocalTime.of(11,57,0), ZoneOffset.UTC);
      if(dateService.now().isAfter(pricesUpdateTime)) {
        cacheString = CACHE_VAR + ", max-age=60";
      } else {
        expiresValue = DateUtils.getGmtStringForCache(pricesUpdateTime);
        cacheString = CACHE_VAR;
      }
    }
    return RestResponse.ResponseBuilder.ok(prices)
      .header(CACHE_CONTROL_HEADER, cacheString)
      .header(EXPIRES_HEADER, expiresValue)
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
