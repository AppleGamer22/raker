import { createFileRoute, stripSearchParams, useNavigate } from "@tanstack/react-router";
import type { MapGeoJSONFeature, MapMouseEvent } from "maplibre-gl";
import { useEffect, useId, useState } from "react";
import z from "zod";

import { PostType, type ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { CardContent } from "@/components/ui/card";
import { Map, MapControls, MapPopup, useMap } from "@/components/ui/map";
import { Progress } from "@/components/ui/progress";
import { useUser } from "@/hooks/user-provider";
import { inPWA } from "@/lib/utils";

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

function MarkersLayer({ histories }: { histories: ScrapeResponse[] }) {
	const { username } = useUser();
	const { exclusive } = Route.useSearch();
	const linkTarget = inPWA() ? undefined : "_blank";
	const { map, isLoaded } = useMap();
	const id = useId();
	const sourceId = `markers-source-${id}`;
	const layerId = `markers-layer-${id}`;
	const [selectedPoint, setSelectedPoint] = useState<ScrapeResponse | null>(null);

	useEffect(() => {
		if (!map || !isLoaded) return;
		map.addSource(sourceId, {
			type: "geojson",
			data: {
				type: "FeatureCollection" as const,
				features: histories.map((history, i) => ({
					type: "Feature" as const,
					properties: {
						i,
					},
					geometry: {
						type: "Point" as const,
						coordinates: [history.coordinates?.longitude ?? 0, history.coordinates?.latitude ?? 0],
					},
				})),
			},
		});

		map.addLayer({
			id: layerId,
			type: "circle",
			source: sourceId,
			paint: {
				"circle-radius": 6,
				"circle-color": "#3b82f6",
				"circle-stroke-width": 2,
				"circle-stroke-color": "#ffffff",
				// add more paint properties here to customize the appearance of the markers
			},
		});

		const handleClick = (
			e: MapMouseEvent & {
				features?: MapGeoJSONFeature[];
			},
		) => {
			if (!e.features?.length) return;

			const feature = e.features[0];
			// const coords = (feature.geometry as GeoJSON.Point).coordinates as [number, number];

			setSelectedPoint(histories[feature.properties.i as number]);
		};

		const handleMouseEnter = () => {
			map.getCanvas().style.cursor = "pointer";
		};

		const handleMouseLeave = () => {
			map.getCanvas().style.cursor = "";
		};

		map.on("click", layerId, handleClick);
		map.on("mouseenter", layerId, handleMouseEnter);
		map.on("mouseleave", layerId, handleMouseLeave);

		return () => {
			map.off("click", layerId, handleClick);
			map.off("mouseenter", layerId, handleMouseEnter);
			map.off("mouseleave", layerId, handleMouseLeave);

			try {
				if (map.getLayer(layerId)) map.removeLayer(layerId);
				if (map.getSource(sourceId)) map.removeSource(sourceId);
			} catch {
				// ignore cleanup errors
			}
		};
	}, [map, isLoaded, sourceId, layerId, histories]);

	return (
		<>
			{selectedPoint && (
				<MapPopup
					latitude={selectedPoint.coordinates?.latitude ?? 0}
					longitude={selectedPoint.coordinates?.longitude ?? 0}
					onClose={() => setSelectedPoint(null)}
					closeOnClick={true}
					focusAfterOpen={false}
					offset={10}
				>
					<HistoryCard
						key={`${selectedPoint}-${selectedPoint.postOwner}-${selectedPoint.post}`}
						history={selectedPoint}
						exclusive={exclusive}
						linkTarget={linkTarget}
						username={username ?? undefined}
					/>
				</MapPopup>
			)}
		</>
	);
}

function MapSearch() {
	const { types, exclusive, categories, owners, page } = Route.useSearch();
	const navigate = useNavigate({ from: Route.fullPath });
	const [histories, setHistories] = useState<ScrapeResponse[]>([]);
	const [_totalCount, setTotalCount] = useState(0n);
	const [isSearching, setIsSearching] = useState(false);
	return (
		<CardContent className="h-[calc(100dvh-2*var(--header-height))] overflow-hidden sm:h-[calc(100dvh-var(--header-height))]">
			<HistorySearchForm
				owners={owners}
				categories={categories}
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
			<Map className="rounded-xl">
				{/* {histories.map((history) => (
					<MapMarker
						key={`${history.postType}-${history.postOwner}-${history.post}`}
						latitude={history.coordinates?.latitude ?? 0}
						longitude={history.coordinates?.longitude ?? 0}
					>
						<MarkerContent>
							<div className="size-4 rounded-full border-2 border-black bg-primary shadow-lg dark:border-white" />
						</MarkerContent>
						<MarkerTooltip>
							<ResultHeader categories={categories} exclusive={exclusive} result={history} showPost />
						</MarkerTooltip>
						<MarkerPopup>
							<HistoryCard
								key={`${history.postType}-${history.postOwner}-${history.post}`}
								history={history}
								exclusive={exclusive}
								linkTarget={linkTarget}
								username={username ?? undefined}
							/>
						</MarkerPopup>
					</MapMarker>
				))} */}
				<MarkersLayer histories={histories} />
				<MapControls showCompass showZoom />
			</Map>
		</CardContent>
	);
}
