import { createFileRoute, stripSearchParams, useNavigate } from "@tanstack/react-router";
import * as maplibregl from "maplibre-gl";

import "maplibre-gl/dist/maplibre-gl.css";
import workerUrl from "maplibre-gl/dist/maplibre-gl-worker.mjs?worker&url";
import { useMemo, useState, useCallback, createElement, type ComponentProps } from "react";
import Map, {
	Popup,
	Source,
	Layer,
	type LayerProps,
	type MapLayerMouseEvent,
	NavigationControl,
} from "react-map-gl/maplibre";
import z from "zod";

maplibregl.setWorkerUrl(new URL(workerUrl, document.baseURI).href);

import { PostType, type ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { CardContent } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { useTheme } from "@/hooks/theme-provider";
import { useUser } from "@/hooks/user-provider";
import { inPWA, primaryColor } from "@/lib/utils";

import { HistoryCard, historySearchDefaults, HistorySearchForm } from "./history";

export const Route = createFileRoute("/map")({
	component: MapSearch,
	validateSearch: z.object({
		types: z.array(z.enum(PostType)).catch(historySearchDefaults.types),
		exclusive: z.boolean().catch(historySearchDefaults.exclusive),
		categories: z.array(z.string()).catch(historySearchDefaults.categories),
		page: z.coerce.bigint().min(1n).catch(historySearchDefaults.page),
		owners: z
			.array(
				z.object({
					owner: z.string(),
					type: z.union([z.enum(PostType), z.literal(-1)]),
				}),
			)
			.catch(historySearchDefaults.owners),
	}),
	search: {
		middlewares: [stripSearchParams(historySearchDefaults)],
	},
});

function CleanSource(props: ComponentProps<typeof Source>) {
	const cleanProps: Record<string, unknown> = {};
	for (const [key, value] of Object.entries(props)) {
		if (!key.startsWith("data-")) {
			cleanProps[key] = value;
		}
	}
	return createElement(Source, cleanProps as ComponentProps<typeof Source>);
}

function CleanLayer(props: ComponentProps<typeof Layer>) {
	const cleanProps: Record<string, unknown> = {};
	for (const [key, value] of Object.entries(props)) {
		if (!key.startsWith("data-")) {
			cleanProps[key] = value;
		}
	}
	return createElement(Layer, cleanProps as ComponentProps<typeof Layer>);
}

function MapSearch() {
	const { types, exclusive, categories, owners, page } = Route.useSearch();
	const navigate = useNavigate({ from: Route.fullPath });
	const { username } = useUser();
	const { computedTheme } = useTheme();
	const linkTarget = inPWA() ? undefined : "_blank";

	const [histories, setHistories] = useState<ScrapeResponse[]>([]);
	const [_totalCount, setTotalCount] = useState(0n);
	const [isSearching, setIsSearching] = useState(false);
	const [selectedPoint, setSelectedPoint] = useState<ScrapeResponse | null>(null);
	const [cursor, setCursor] = useState<string>("");

	const geojsonData = useMemo(() => {
		return {
			type: "FeatureCollection" as const,
			features: histories.map((history, i) => ({
				type: "Feature" as const,
				properties: { i },
				geometry: {
					type: "Point" as const,
					coordinates: [history.coordinates?.longitude ?? 0, history.coordinates?.latitude ?? 0],
				},
			})),
		};
	}, [histories]);

	const layerStyle = useMemo<LayerProps>(() => {
		return {
			id: "markers-layer",
			type: "circle",
			paint: {
				"circle-radius": 6,
				"circle-color": primaryColor,
				"circle-stroke-width": 2,
				"circle-stroke-color": computedTheme === "dark" ? "white" : "black",
			},
		};
	}, [computedTheme]);

	const mapStyle =
		computedTheme === "dark"
			? "https://basemaps.cartocdn.com/gl/dark-matter-gl-style/style.json"
			: "https://basemaps.cartocdn.com/gl/positron-gl-style/style.json";

	const onClick = useCallback(
		(event: MapLayerMouseEvent) => {
			const feature = event.features?.[0];
			if (feature) {
				const historyIndex = feature.properties?.i as number;
				if (historyIndex !== undefined && histories[historyIndex]) {
					setSelectedPoint(histories[historyIndex]);
				}
			} else {
				setSelectedPoint(null);
			}
		},
		[histories],
	);

	const onMouseEnter = useCallback(() => setCursor("pointer"), []);
	const onMouseLeave = useCallback(() => setCursor(""), []);

	return (
		<CardContent className="flex h-[calc(100dvh-2*var(--header-height))] flex-col overflow-hidden sm:h-[calc(100dvh-var(--header-height))]">
			<HistorySearchForm
				owners={owners}
				categories={categories}
				availablePostTypes={[PostType.VSCO]}
				exclusive={exclusive}
				types={types}
				currentPage={page}
				setHistories={setHistories}
				setTotalCount={setTotalCount}
				setHistorySearchPending={setIsSearching}
				onSearchSubmit={async ({ categories, exclusive, ownersSearchValue, types, page }) => {
					await navigate({
						search: {
							types,
							exclusive,
							categories,
							page,
							owners: ownersSearchValue,
						},
						replace: true,
					});
				}}
				autoSubmit={false}
				onlyWithCoordinates
			/>
			{isSearching && <Progress className="pt-2" value={null} />}
			<div className="relative isolate min-h-0 flex-1 overflow-hidden rounded-xl">
				<Map
					mapLib={maplibregl}
					initialViewState={{
						longitude: 0,
						latitude: 0,
						zoom: 2,
					}}
					mapStyle={mapStyle}
					interactiveLayerIds={["markers-layer"]}
					onClick={onClick}
					onMouseEnter={onMouseEnter}
					onMouseLeave={onMouseLeave}
					cursor={cursor}
					attributionControl={false}
					renderWorldCopies={false}
				>
					<CleanSource id="markers-source" type="geojson" data={geojsonData}>
						<CleanLayer {...layerStyle} />
					</CleanSource>

					<NavigationControl position="bottom-right" />

					{selectedPoint && (
						<Popup
							longitude={selectedPoint.coordinates?.longitude ?? 0}
							latitude={selectedPoint.coordinates?.latitude ?? 0}
							onClose={() => setSelectedPoint(null)}
							closeOnClick={true}
							closeButton={false}
							focusAfterOpen={false}
							offset={10}
							maxWidth="none"
							className="[&_.maplibregl-popup-content]:bg-transparent [&_.maplibregl-popup-content]:p-0 [&_.maplibregl-popup-content]:shadow-none [&_.maplibregl-popup-tip]:hidden"
						>
							<div className="relative max-w-62 animate-in rounded-md border bg-popover p-3 text-popover-foreground shadow-md duration-200 ease-out fade-in-0 zoom-in-95">
								<HistoryCard
									key={`${selectedPoint.postType}-${selectedPoint.postOwner}-${selectedPoint.post}`}
									history={selectedPoint}
									exclusive={exclusive}
									linkTarget={linkTarget}
									username={username ?? undefined}
								/>
							</div>
						</Popup>
					)}
				</Map>
			</div>
		</CardContent>
	);
}
