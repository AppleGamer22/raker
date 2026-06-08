import { createFileRoute } from "@tanstack/react-router";

import { CardContent } from "@/components/ui/card";
import { Map, MapControls } from "@/components/ui/map";

export const Route = createFileRoute("/map")({
	component: MapSearch,
});

function MapSearch() {
	return (
		<CardContent className="h-[calc(100dvh-2*var(--header-height))] overflow-hidden sm:h-[calc(100dvh-var(--header-height))]">
			<Map className="rounded-xl">
				<MapControls showCompass showLocate showZoom />
			</Map>
		</CardContent>
	);
}
