package fi.ehin.utils;

public final class RequestUtils {
    private RequestUtils() {}

    public static final String CACHE_CONTROL_HEADER = "Cache-Control";
    public static final String EXPIRES_HEADER = "Expires";


    public static final String CACHE_VAR = "public";
    public static final String CACHE_LONG =
            "public, max-age=" + 60 * 60 * 168 + ", immutable";
}
