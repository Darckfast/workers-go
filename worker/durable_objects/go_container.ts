import { Container } from "@cloudflare/containers";

// docker build -t cloudflare-dev/gocontainer ./worker
export class GoContainer extends Container {
  defaultPort = 8080; // The default port for the container to listen on
  sleepAfter = '6m'; // Sleep the container if no requests are made in this timeframe

  override onStart() {
    console.log('Container successfully started');
  }

  override onStop() {
    console.log('Container successfully shut down');
  }

  override onError(error: unknown) {
    console.log('Container error:', error);
  }
}
