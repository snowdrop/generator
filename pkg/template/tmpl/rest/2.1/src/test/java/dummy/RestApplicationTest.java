/*
 * Copyright 2016-2017 Red Hat, Inc, and individual contributors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package {{.PackageName}};

import static io.restassured.RestAssured.given;
import static org.hamcrest.core.Is.is;

import {{.PackageName}}.service.Greeting;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;


@RunWith(SpringRunner.class)
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
public abstract class RestApplicationTest {
    @Value("${local.server.port}")
    private int port;

    private static final String GREETING_PATH = "api/greeting";

    @Test
    public void testGreetingEndpoint() {
        given()
           .baseUri(baseURI())
           .get(GREETING_PATH)
           .then()
           .statusCode(200)
           .body("content", is(String.format(Greeting.FORMAT, "World")));
    }

    @Test
    public void testGreetingEndpointWithNameParameter() {
        given()
           .baseUri(baseURI())
           .param("name", "John")
           .when()
           .get(GREETING_PATH)
           .then()
           .statusCode(200)
           .body("content", is(String.format(Greeting.FORMAT, "John")));
    }

    private String baseURI() {
        return String.format("http://localhost:%d", port);
    }
}
