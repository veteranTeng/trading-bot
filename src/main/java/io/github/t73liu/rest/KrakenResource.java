package io.github.t73liu.rest;

import io.github.t73liu.model.ExceptionWrapper;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiResponse;
import io.swagger.annotations.ApiResponses;
import org.glassfish.jersey.client.ClientConfig;
import org.glassfish.jersey.jackson.JacksonFeature;
import org.springframework.stereotype.Component;

import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.client.Client;
import javax.ws.rs.client.ClientBuilder;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import java.util.Map;

@Component
@Path("/kraken")
@Consumes(MediaType.APPLICATION_JSON)
@Produces(MediaType.APPLICATION_JSON)
@Api("KrakenResource")
@ApiResponses(@ApiResponse(code = 500, message = "Internal Server Error", response = ExceptionWrapper.class))
public class KrakenResource {
    @GET
    @Path("/test")
    @ApiResponses(@ApiResponse(code = 200, message = "Retrieved Currency Information", responseContainer = "List", response = Map.class))
    public Response test() {
        ClientConfig cc = new ClientConfig().register(new JacksonFeature());
        Client client = ClientBuilder.newClient(cc);
        return Response.ok(client.target("https://api.kraken.com/0/public/AssetPairs").request().get().getEntity()).build();
    }
}
