package io.github.yufanfan29.ripple.content;

import static org.assertj.core.api.Assertions.assertThat;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.availability.AvailabilityChangeEvent;
import org.springframework.boot.availability.ReadinessState;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment;
import org.springframework.context.ApplicationContext;

@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
class HealthCheckTests {

  private final HttpClient client = HttpClient.newHttpClient();
  private final ApplicationContext context;
  private final int port;

  HealthCheckTests(ApplicationContext context, @Value("${local.server.port}") int port) {
    this.context = context;
    this.port = port;
  }

  @AfterEach
  void acceptTrafficAgain() {
    // The tests share one application context, so reset readiness for the next test.
    AvailabilityChangeEvent.publish(context, ReadinessState.ACCEPTING_TRAFFIC);
  }

  @Test
  void livezReportsUp() throws IOException, InterruptedException {
    assertThat(get("/livez").statusCode()).isEqualTo(200);
  }

  @Test
  void readyzReportsUp() throws IOException, InterruptedException {
    assertThat(get("/readyz").statusCode()).isEqualTo(200);
  }

  @Test
  void readyzIsUnavailableWhileRefusingTraffic() throws IOException, InterruptedException {
    refuseTraffic();

    assertThat(get("/readyz").statusCode()).isEqualTo(503);
  }

  @Test
  void livezStaysUpWhileRefusingTraffic() throws IOException, InterruptedException {
    refuseTraffic();

    assertThat(get("/livez").statusCode()).isEqualTo(200);
  }

  // Spring Boot publishes the same event when shutdown starts.
  private void refuseTraffic() {
    AvailabilityChangeEvent.publish(context, ReadinessState.REFUSING_TRAFFIC);
  }

  private HttpResponse<String> get(String path) throws IOException, InterruptedException {
    HttpRequest request =
        HttpRequest.newBuilder(URI.create("http://localhost:" + port + path)).build();
    return client.send(request, HttpResponse.BodyHandlers.ofString());
  }
}
