package fi.ehin.resource;

import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;

@Path("/hello")
public class GreetingResource {

  @GET
  @Produces(MediaType.TEXT_PLAIN)
  public String hello() {
    printMemory();
    return "Hello from Quarkus REST";
  }

  public void printMemory() {
    System.out.println(
      "Meg used=" +
      (Runtime.getRuntime().totalMemory() - Runtime.getRuntime().freeMemory()) /
      (1000 * 1000) +
      "M"
    );
  }
}
