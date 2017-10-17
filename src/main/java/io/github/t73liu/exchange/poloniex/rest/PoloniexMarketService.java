package io.github.t73liu.exchange.poloniex.rest;

import eu.verdelhan.ta4j.BaseTick;
import eu.verdelhan.ta4j.Tick;
import io.github.t73liu.exchange.ExchangeService;
import io.github.t73liu.exchange.MarketService;
import io.github.t73liu.model.CandlestickIntervals;
import io.github.t73liu.model.poloniex.PoloniexCandle;
import io.github.t73liu.model.poloniex.PoloniexPair;
import io.github.t73liu.util.DateUtil;
import it.unimi.dsi.fastutil.objects.ObjectArrayList;
import org.apache.http.NameValuePair;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.message.BasicNameValuePair;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.Arrays;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import static io.github.t73liu.util.HttpUtil.generateGet;
import static io.github.t73liu.util.MapperUtil.JSON_READER;

@Service
@ConfigurationProperties(prefix = "poloniex")
public class PoloniexMarketService extends ExchangeService implements MarketService {
    public static Tick mapExchangeCandleToTick(PoloniexCandle candle) {
        return new BaseTick(DateUtil.unixTimeStampToZonedDateTime(candle.getDate()), candle.getOpen(), candle.getHigh(), candle.getLow(), candle.getClose(), candle.getVolume());
    }

    public Map<String, Map<String, String>> getAllTicker() throws Exception {
        List<NameValuePair> queryParams = new ObjectArrayList<>(1);
        queryParams.add(new BasicNameValuePair("command", "returnTicker"));
        HttpGet get = generateGet(getPublicUrl(), queryParams);

        try (CloseableHttpClient httpClient = HttpClients.createDefault();
             CloseableHttpResponse response = httpClient.execute(get)) {
            return JSON_READER.readValue(response.getEntity().getContent());
        } finally {
            get.releaseConnection();
        }
    }

    public Map getTickerForPair(PoloniexPair pair) throws Exception {
        return getAllTicker().get(pair.getPairName());
    }

    public Map getOrderBookForPair(PoloniexPair pair, int amount) throws Exception {
        List<NameValuePair> queryParams = new ObjectArrayList<>(3);
        queryParams.add(new BasicNameValuePair("command", "returnOrderBook"));
        queryParams.add(new BasicNameValuePair("depth", String.valueOf(amount)));
        queryParams.add(new BasicNameValuePair("currencyPair", pair.getPairName()));
        HttpGet get = generateGet(getPublicUrl(), queryParams);

        try (CloseableHttpClient httpClient = HttpClients.createDefault();
             CloseableHttpResponse response = httpClient.execute(get)) {
            return JSON_READER.readValue(response.getEntity().getContent());
        } finally {
            get.releaseConnection();
        }
    }

    public Map getAllOrderBook(int amount) throws Exception {
        List<NameValuePair> queryParams = new ObjectArrayList<>(3);
        queryParams.add(new BasicNameValuePair("command", "returnOrderBook"));
        queryParams.add(new BasicNameValuePair("depth", String.valueOf(amount)));
        queryParams.add(new BasicNameValuePair("currencyPair", "all"));
        HttpGet get = generateGet(getPublicUrl(), queryParams);

        try (CloseableHttpClient httpClient = HttpClients.createDefault();
             CloseableHttpResponse response = httpClient.execute(get)) {
            return JSON_READER.readValue(response.getEntity().getContent());
        } finally {
            get.releaseConnection();
        }
    }

    public List<Tick> getCandlestick(PoloniexPair pair, LocalDateTime startDateTime, LocalDateTime endDateTime, CandlestickIntervals period) throws Exception {
        return Arrays.stream(getExchangeCandle(pair, startDateTime, endDateTime, period))
                .map(PoloniexMarketService::mapExchangeCandleToTick)
                .collect(Collectors.toList());
    }

    public PoloniexCandle[] getExchangeCandle(PoloniexPair pair, LocalDateTime startDateTime, LocalDateTime endDateTime, CandlestickIntervals period) throws Exception {
        List<NameValuePair> queryParams = new ObjectArrayList<>(3);
        queryParams.add(new BasicNameValuePair("command", "returnChartData"));
        // Candlestick period in seconds 300,900,1800,7200,14400,86400
        queryParams.add(new BasicNameValuePair("period", String.valueOf(period.getInterval())));
        queryParams.add(new BasicNameValuePair("currencyPair", pair.getPairName()));
        // UNIX timestamp format of specified time range (i.e. last hour)
        queryParams.add(new BasicNameValuePair("start", String.valueOf(DateUtil.localDateTimeToUnixTimestamp(startDateTime))));
        queryParams.add(new BasicNameValuePair("end", String.valueOf(DateUtil.localDateTimeToUnixTimestamp(endDateTime))));
        HttpGet get = generateGet(getPublicUrl(), queryParams);

        try (CloseableHttpClient httpClient = HttpClients.createDefault();
             CloseableHttpResponse response = httpClient.execute(get)) {
            return JSON_READER.forType(PoloniexCandle[].class).readValue(response.getEntity().getContent());
        } finally {
            get.releaseConnection();
        }
    }

    public Map getCurrencies() throws Exception {
        List<NameValuePair> queryParams = new ObjectArrayList<>(1);
        queryParams.add(new BasicNameValuePair("command", "returnCurrencies"));
        HttpGet get = generateGet(getPublicUrl(), queryParams);

        try (CloseableHttpClient httpClient = HttpClients.createDefault();
             CloseableHttpResponse response = httpClient.execute(get)) {
            return JSON_READER.readValue(response.getEntity().getContent());
        } finally {
            get.releaseConnection();
        }
    }

    public Map getLoanOrders(String currency) throws Exception {
        List<NameValuePair> queryParams = new ObjectArrayList<>(2);
        queryParams.add(new BasicNameValuePair("command", "returnLoanOrders"));
        queryParams.add(new BasicNameValuePair("currency", currency));
        HttpGet get = generateGet(getPublicUrl(), queryParams);

        try (CloseableHttpClient httpClient = HttpClients.createDefault();
             CloseableHttpResponse response = httpClient.execute(get)) {
            return JSON_READER.readValue(response.getEntity().getContent());
        } finally {
            get.releaseConnection();
        }
    }

    private String getPublicUrl() {
        return getBaseUrl() + "public";
    }
}
