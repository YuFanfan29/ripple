package io.github.yufanfan29.ripple.content;

import static org.assertj.core.api.Assertions.assertThat;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment;

@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
class HealthCheckTests {

  private final HttpClient client = HttpClient.newHttpClient();
  private final int port;

  HealthCheckTests(@Value("${local.server.port}") int port) {
    this.port = port;
  }

  @Test
  void livezReportsUp() throws IOException, InterruptedException {
    assertThat(get("/livez").statusCode()).isEqualTo(200);
  }

  @Test
  void readyzReportsUp() throws IOException, InterruptedException {
    assertThat(get("/readyz").statusCode()).isEqualTo(200);
  }

  private HttpResponse<String> get(String path) throws IOException, InterruptedException {
    HttpRequest request =
        HttpRequest.newBuilder(URI.create("http://localhost:" + port + path)).build();
    return client.send(request, HttpResponse.BodyHandlers.ofString());
  }
}
