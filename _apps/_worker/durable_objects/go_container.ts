import { Container } from "@cloudflare/containers";
/**
 * Since Cloudflare Container's is new and going to constant changes
 * its recommended to follow the official docs https://developers.cloudflare.com/containers
 */
export class GoContainer extends Container {
	defaultPort = 8080; // The default port for the container to listen on
	sleepAfter = "6m"; // Sleep the container if no requests are made in this timeframe

	override onStart() {
		console.log("Container successfully started");
	}

	override onStop() {
		console.log("Container successfully shut down");
	}

	override onError(error: unknown) {
		console.log("Container error:", error);
	}
}
