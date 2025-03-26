package fi.ehin.resource;

import fi.ehin.repository.PriceRepository;
import fi.ehin.service.PricesService;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.PathParam;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.QueryParam;
import jakarta.ws.rs.core.MediaType;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.OffsetDateTime;
import java.time.OffsetTime;
import java.time.ZoneId;
import java.util.List;
import org.eclipse.microprofile.config.inject.ConfigProperty;
import org.jboss.resteasy.reactive.RestResponse;

@Path("/api")
public class PriceResource {

  private final PriceRepository priceRepository;
  private final PricesService pricesService;

  private static final String CACHE_ONE_MINUTE =
    "public, max-age=" + 60 + ", immutable";
  private static final String CACHE_TEN_HOUR =
    "public, max-age=" + 60 * 60 * 10 + ", immutable";

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
    final var zone = ZoneId.of("Europe/Helsinki");
    final var dateWithTime = date.atTime(
      OffsetTime.of(0, 0, 0, 0, zone.getRules().getOffset(LocalDateTime.now()))
    );

    final var prices = priceRepository.getPrices(
      dateWithTime.minusDays(1),
      dateWithTime.plusDays(2)
    );
    String cacheString = CACHE_TEN_HOUR; // 10 hours cache if already have days prices
    final var newestPriceDate = prices.getLast().deliveryStart().toLocalDate();
    if (!newestPriceDate.equals(date.plusDays(1))) {
      // No tomorrow's prices yet so cache only one minute
      cacheString = CACHE_ONE_MINUTE;
    }
    return RestResponse.ResponseBuilder.ok(prices)
      .header("Cache-Control", cacheString)
      .build();
  }

  public record UpdatePricesResponse(OffsetDateTime time, boolean done) {}

  @GET
  @Produces(MediaType.APPLICATION_JSON)
  @Path("/update-prices")
  public RestResponse<UpdatePricesResponse> updatePrices(
    @QueryParam("p") final String password
  ) {
    if (!updatePricesPassword.equals(password)) {
      return RestResponse.notFound();
    }
    final var zone = ZoneId.of("Europe/Helsinki");
    final var dateWithTime = LocalDate.now()
      .atTime(
        OffsetTime.of(
          0,
          0,
          0,
          0,
          zone.getRules().getOffset(LocalDateTime.now())
        )
      );
    final var prices = pricesService.GetTomorrowsPrices();
    if (prices == null) {
      return RestResponse.ok(new UpdatePricesResponse(dateWithTime, false));
    }
    priceRepository.insertPrices(prices.multiAreaEntries());

    return RestResponse.ok(new UpdatePricesResponse(dateWithTime, true));
  }

  @POST
  @Produces(MediaType.APPLICATION_JSON)
  @Path("/update-prices/{date}")
  public RestResponse<UpdatePricesResponse> updatePricesForDate(
    @PathParam("date") final LocalDate date,
    @QueryParam("p") final String password
  ) {
    if (!updatePricesPassword.equals(password)) {
      return RestResponse.notFound();
    }
    final var zone = ZoneId.of("Europe/Helsinki");
    final var dateWithTime = LocalDate.now()
      .atTime(
        OffsetTime.of(
          0,
          0,
          0,
          0,
          zone.getRules().getOffset(LocalDateTime.now())
        )
      );
    final var prices = pricesService.GetPrices(date);
    priceRepository.insertPrices(prices.multiAreaEntries());

    return RestResponse.ok(new UpdatePricesResponse(dateWithTime, true));
  }
}
