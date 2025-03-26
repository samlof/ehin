package fi.ehin.service;

import static org.hamcrest.CoreMatchers.*;
import static org.hamcrest.MatcherAssert.assertThat;

import io.quarkus.test.junit.QuarkusTest;
import jakarta.inject.Inject;
import java.time.LocalDate;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

@QuarkusTest
public class PricesServiceTest {

  @Inject
  PricesService pricesService;

  @Test
  void testGettingPrices() {
    final var prices = pricesService.GetTomorrowsPrices();
    Assertions.assertEquals("DayAhead", prices.market());
    Assertions.assertEquals(
      LocalDate.now().plusDays(1),
      prices.deliveryDateCET()
    );
    Assertions.assertEquals("EUR", prices.currency());
    assertThat(prices.deliveryAreas(), hasItem("FI"));
  }
}
