import { createFileRoute, stripSearchParams, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import z from "zod";

import { PostType, type ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { CardContent } from "@/components/ui/card";
import { Map, MapControls } from "@/components/ui/map";

import { historySearchDefaults, HistorySearchForm } from "./history";
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

function MapSearch() {
	const { types, exclusive, categories, owners, page } = Route.useSearch();
	const navigate = useNavigate({ from: Route.fullPath });
	const [_histories, setHistories] = useState<ScrapeResponse[]>([]);
	const [_totalCount, setTotalCount] = useState(0n);
	const [_isSearching, setIsSearching] = useState(false);
	return (
		<CardContent className="h-[calc(100dvh-2*var(--header-height))] overflow-hidden sm:h-[calc(100dvh-var(--header-height))]">
			<HistorySearchForm
				owners={owners}
				categories={categories}
				exclusive={exclusive}
				types={types}
				currentPage={page}
				autoSubmit={true}
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
			/>
			<Map className="rounded-xl">
				<MapControls showCompass showZoom />
			</Map>
		</CardContent>
	);
}
