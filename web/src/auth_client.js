
import { authEvents, LOGOUT_EVENT } from './auth_event';

// State to track token refreshing
let isRefreshing = false;
let failedQueue = [];

// Helper to process the queued requests once the token is refreshed
const processQueue = (error, token = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

const getAccessToken = () => sessionStorage.getItem("accessToken");
const getRefreshToken = () => sessionStorage.getItem("refreshToken");
const updateTokens = (access, refresh) => {
  sessionStorage.setItem("accessToken", access);
  sessionStorage.setItem("refreshToken", refresh);
};
const clearTokens = () => {
  sessionStorage.removeItem("accessToken");
  sessionStorage.removeItem("refreshToken");
};


// The core refresh token network call
async function handleTokenRefresh() {
  const refreshToken = getRefreshToken();
  if (!refreshToken) throw new Error("No refresh token available");

  const response = await fetch("http://localhost:8080/refresh/", {
    mode: 'cors',
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refreshToken }),
  });

  if (!response.ok) {
    throw new Error("Refresh token expired or invalid");
  }

  const data = await response.json();
  return data; // Assuming returns { accessToken, refreshToken }
}


// Custom interceptor wrapper around native fetch
export async function fetchClient(url, options = {}) {
  // Ensure headers object exists
  options.headers = options.headers || {};

  // Request Interceptor: Inject the current access token
  const token = getAccessToken();
  if (token) {
    options.headers["Authorization"] = `Bearer ${token}`;
  }

  try {
    const response = await fetch(url, options);

    // Response Interceptor: Catch 401 Unauthorized status
    if (response.status === 401) {
      
      // If a refresh is already in progress, queue this request
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((newAccessToken) => {
            options.headers["Authorization"] = `Bearer ${newAccessToken}`;
            return fetch(url, options); // Retry with new token
          })
          .catch((err) => Promise.reject(err));
      }

      // Mark that refreshing has started to throttle other incoming 401s
      isRefreshing = true;

      return new Promise((resolve, reject) => {
        handleTokenRefresh().then((data) => {
            updateTokens(data.accessToken, sessionStorage.getItem("refreshToken")); // data.refreshToken);
            
            // Update original request header and retry it
            options.headers["Authorization"] = `Bearer ${data.accessToken}`;
            
            processQueue(null, data.accessToken);
            
            resolve(fetch(url, options));

          }).catch((refreshError) => {          
            // refreshing the accessToken has failed
            processQueue(refreshError, null);            
            
            // remove the tokens from sessionStorage
            clearTokens();
            
            // dispatch LOGOUT_EVENT to pass control 
            // back to view to call setLoggedOutState
            authEvents.dispatchEvent(new Event(LOGOUT_EVENT));
            
            reject(refreshError);
          })
          .finally(() => {
            isRefreshing = false;
          });
      });
    }

    return response;
  } catch (networkError) {
    return Promise.reject(networkError);
  }
}
