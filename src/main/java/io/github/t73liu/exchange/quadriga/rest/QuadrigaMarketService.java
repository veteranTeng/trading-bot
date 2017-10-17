package io.github.t73liu.exchange.quadriga.rest;

import io.github.t73liu.exchange.ExchangeService;
import io.github.t73liu.exchange.MarketService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Service;

@Service
@ConfigurationProperties(prefix = "quadrigacx")
public class QuadrigaMarketService extends ExchangeService implements MarketService {
    private static final Logger LOGGER = LoggerFactory.getLogger(QuadrigaMarketService.class);
    private static final double FIAT_FEE = 0.005;
    private static final double CRYPTO_FEE = 0.002;

    // TODO implement
}
