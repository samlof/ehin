package fi.ehin.service;

import jakarta.enterprise.context.ApplicationScoped;
import java.time.OffsetDateTime;

@ApplicationScoped
public class DateService {

  public OffsetDateTime now() {
    return OffsetDateTime.now();
  }
}
