package io.github.yufanfan29.ripple.content;

import static org.assertj.core.api.Assertions.assertThat;

import java.util.Map;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.json.JsonParserFactory;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.system.CapturedOutput;
import org.springframework.boot.test.system.OutputCaptureExtension;

@SpringBootTest
@ExtendWith(OutputCaptureExtension.class)
class LoggingTests {

  private static final Logger log = LoggerFactory.getLogger(LoggingTests.class);

  @Test
  void logLinesAreJsonWithConventionFields(CapturedOutput output) {
    log.info("slf4j test marker");

    Map<String, Object> fields = logLine(output, "slf4j test marker");

    assertThat(fields)
        .containsEntry("level", "INFO")
        .containsEntry("service", "content")
        .containsEntry("message", "slf4j test marker")
        .containsKey("time")
        .doesNotContainKeys("@timestamp", "@version", "level_value");
  }

  @Test
  void frameworkLogLinesHaveNoLoggingApiTags(CapturedOutput output) {
    // Spring and Tomcat log through Commons Logging, which tags every line.
    Log commonsLog = LogFactory.getLog(LoggingTests.class);
    commonsLog.info("commons logging test marker");

    assertThat(logLine(output, "commons logging test marker")).doesNotContainKey("tags");
  }

  private static Map<String, Object> logLine(CapturedOutput output, String marker) {
    String line =
        output
            .getOut()
            .lines()
            .filter(candidate -> candidate.contains(marker))
            .findFirst()
            .orElseThrow();
    return JsonParserFactory.getJsonParser().parseMap(line);
  }
}
