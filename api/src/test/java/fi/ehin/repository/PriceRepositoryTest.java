package fi.ehin.repository;

import static org.hamcrest.CoreMatchers.hasItem;
import static org.hamcrest.MatcherAssert.assertThat;

import fi.ehin.external.NordPoolClient;
import io.quarkus.test.junit.QuarkusTest;
import jakarta.inject.Inject;
import java.math.BigDecimal;
import java.time.OffsetDateTime;
import java.util.List;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

@QuarkusTest
public class PriceRepositoryTest {

  @Inject
  PriceRepository priceRepository;

  @Test
  void testSelect1Works() {
    final var one = priceRepository.select1();
    Assertions.assertEquals(1, one);
  }
}
