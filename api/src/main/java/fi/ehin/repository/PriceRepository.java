package fi.ehin.repository;

import com.fasterxml.jackson.annotation.JsonProperty;
import fi.ehin.external.NordPoolClient;
import io.agroal.api.AgroalDataSource;
import io.quarkus.logging.Log;
import io.quarkus.runtime.annotations.RegisterForReflection;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.transaction.Transactional;
import java.math.BigDecimal;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.time.OffsetDateTime;
import java.util.ArrayList;
import java.util.List;

@RegisterForReflection(
  targets = { org.postgresql.ssl.DefaultJavaSSLFactory.class }
)
@ApplicationScoped
public class PriceRepository {

  private final AgroalDataSource dataSource;

  public PriceRepository(AgroalDataSource dataSource) {
    this.dataSource = dataSource;
  }

  public int select1() {
    try (
      final var con = dataSource.getConnection();
      final var ps = con.prepareStatement(
        """
        SELECT 1
        """
      );
      final var rs = ps.executeQuery()
    ) {
      rs.next();
      return rs.getInt(1);
    } catch (SQLException e) {
      throw new RuntimeException(e);
    }
  }

  public record PriceHistoryEntry(
    @JsonProperty("p") BigDecimal price,
    @JsonProperty("s") OffsetDateTime deliveryStart,
    @JsonProperty("e") OffsetDateTime deliveryEnd
  ) {
    public static PriceHistoryEntry fromSql(final ResultSet rs)
      throws SQLException {
      return new PriceHistoryEntry(
        rs.getObject("price", BigDecimal.class),
        rs.getObject("delivery_start", OffsetDateTime.class),
        rs.getObject("delivery_end", OffsetDateTime.class)
      );
    }
  }

  public List<PriceHistoryEntry> getPrices(
    final OffsetDateTime from,
    final OffsetDateTime to
  ) {
    try (
      final var con = dataSource.getConnection();
      final var ps = con.prepareStatement(
        """
        SELECT price, delivery_start, delivery_end
        FROM price_history
        WHERE delivery_start > ? AND delivery_end < ?
        ORDER BY delivery_start
        """
      )
    ) {
      ps.setObject(1, from);
      ps.setObject(2, to);
      try (final var rs = ps.executeQuery()) {
        final var ret = new ArrayList<PriceHistoryEntry>();
        while (rs.next()) {
          ret.add(PriceHistoryEntry.fromSql(rs));
        }
        return ret;
      }
    } catch (SQLException e) {
      throw new RuntimeException(e);
    }
  }

  @Transactional
  public void insertPrices(
    List<NordPoolClient.PriceDataResponseEntry> entryList
  ) {
    try (final var con = dataSource.getConnection()) {
      final var sb = new StringBuilder();
      sb.append(
        "INSERT INTO price_history (price, delivery_start, delivery_end)\nVALUES\n"
      );
      for (int i = 0; i < entryList.size(); i++) {
        sb.append("(?, ?, ?)");
        if (i < entryList.size() - 1) {
          sb.append(", \n");
        }
      }
      sb.append(" ON CONFLICT (delivery_start) DO NOTHING");
      try (final var ps = con.prepareStatement(sb.toString())) {
        int paramCounter = 1;
        for (final var entry : entryList) {
          ps.setBigDecimal(paramCounter++, entry.entryPerArea().FI());
          ps.setObject(paramCounter++, entry.deliveryStart());
          ps.setObject(paramCounter++, entry.deliveryEnd());
        }

        ps.execute();
        final var rowsAffected = ps.getUpdateCount();
        Log.infof(
          "Added %d of %d prices to db",
          rowsAffected,
          entryList.size()
        );
      }
    } catch (SQLException e) {
      throw new RuntimeException(e);
    }
  }
}
