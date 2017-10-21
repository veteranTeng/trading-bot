package io.github.t73liu.exchange.bitfinex.rest;

import io.github.t73liu.exchange.ExchangeService;
import io.github.t73liu.exchange.MarketService;
import io.github.t73liu.model.bitfinex.BitfinexCandle;
import io.github.t73liu.model.bitfinex.BitfinexCandleInterval;
import io.github.t73liu.model.bitfinex.BitfinexPair;
import io.github.t73liu.model.bitfinex.BitfinexTicker;
import it.unimi.dsi.fastutil.objects.ObjectArrayList;
import org.apache.http.NameValuePair;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.message.BasicNameValuePair;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Service;

import java.util.Iterator;
import java.util.List;

import static io.github.t73liu.util.HttpUtil.generateGet;
import static io.github.t73liu.util.MapperUtil.JSON_READER;

@Service
@ConfigurationProperties(prefix = "bitfinex")
public class BitfinexMarketService extends ExchangeService implements MarketService {
    public BitfinexTicker getTickerForPair(BitfinexPair pair) throws Exception {
        return new BitfinexTicker(getExchangeTickerForPair(pair));
    }

    private static List<NameValuePair> getCandleParameters(Long startMilliseconds, Long endMilliseconds, int limit, boolean newestFirst) {
        List<NameValuePair> queryParams = new ObjectArrayList<>(4);
        queryParams.add(new BasicNameValuePair("limit", String.valueOf(limit)));
        if (startMilliseconds != null) {
            queryParams.add(new BasicNameValuePair("start", startMilliseconds.toString()));
        }
        if (endMilliseconds != null) {
            queryParams.add(new BasicNameValuePair("end", endMilliseconds.toString()));
        }
        queryParams.add(new BasicNameValuePair("sort", newestFirst ? "-1" : "1"));
        return queryParams;
    }

    private double[] getExchangeTickerForPair(BitfinexPair pair) throws Exception {
        List<NameValuePair> queryParams = new ObjectArrayList<>(0);
        HttpGet get = generateGet(getTickerUrl(pair.getPairName()), queryParams);

        try (CloseableHttpClient httpClient = HttpClients.createDefault();
             CloseableHttpResponse response = httpClient.execute(get)) {
            return JSON_READER.forType(double[].class).readValue(response.getEntity().getContent());
        } finally {
            get.releaseConnection();
        }
    }

    private String getTickerUrl(String pairName) {
        return getBaseUrl() + "ticker/" + pairName;
    }

    public List<BitfinexCandle> getExchangeCandleForPair(BitfinexPair pair, BitfinexCandleInterval period, Long startMilliseconds,
                                                         Long endMilliseconds, int limit, boolean newestFirst) throws Exception {
        HttpGet get = generateGet(getCandleUrl(pair.getPairName(), period.getIntervalName()),
                getCandleParameters(startMilliseconds, endMilliseconds, limit, newestFirst));

        try (CloseableHttpClient httpClient = HttpClients.createDefault();
             CloseableHttpResponse response = httpClient.execute(get)) {
            List<BitfinexCandle> candles = new ObjectArrayList<>(limit);
            Iterator<double[]> content = JSON_READER.forType(double[].class).readValues(response.getEntity().getContent());
            while (content.hasNext()) {
                candles.add(new BitfinexCandle(content.next()));
            }
            return candles;
        } finally {
            get.releaseConnection();
        }
    }

    private String getCandleUrl(String pairName, String candleInterval) {
        return getBaseUrl() + "candles/trade:" + candleInterval + ":" + pairName + "/hist";
    }
}
