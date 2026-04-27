import { useQuery } from "@connectrpc/connect-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect } from "react";

import { getUserCategories } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { FieldGroup, FieldLegend, Field, FieldSet, FieldLabel } from "@/components/ui/field";
import { InstagramIcon } from "@/components/ui/svgs/instagram";
import { SnapchatIcon } from "@/components/ui/svgs/snapchat";
import { TikTokIcon } from "@/components/ui/svgs/tiktok";
import { VSCOIcon } from "@/components/ui/svgs/vsco";
import { useUser } from "@/hooks/user-provider";

export const Route = createFileRoute("/history")({
	component: History,
});

function HistorySearchForm() {
	// const categoriesQuery = useQuery(getUserCategories, {});
	// const categories = categoriesQuery.data?.categories ?? [];

	return (
		<form>
			<FieldGroup>
				<FieldSet>
					<FieldLegend>Post Types</FieldLegend>
					<FieldGroup
						data-slot="checkbox-group"
						className="flex flex-row flex-wrap gap-1 [&>[data-slot=field]]:w-auto"
					>
						<Field orientation="horizontal">
							<Checkbox id="post-type-instagram" />
							<FieldLabel htmlFor="post-type-instagram">
								<InstagramIcon className="w-4" />
								Instagram Post
							</FieldLabel>
						</Field>
						<Field orientation="horizontal">
							<Checkbox id="post-type-highlight" />
							<FieldLabel htmlFor="post-type-highlight">
								<InstagramIcon className="w-4" />
								Instagram Highlight
							</FieldLabel>
						</Field>
						<Field orientation="horizontal">
							<Checkbox id="post-type-story" />
							<FieldLabel htmlFor="post-type-story">
								<InstagramIcon className="w-4" />
								Instagram Post
							</FieldLabel>
						</Field>
						<Field orientation="horizontal">
							<Checkbox id="post-type-tiktok" />
							<FieldLabel htmlFor="post-type-tiktok">
								<TikTokIcon className="w-4" />
								TikTok Post
							</FieldLabel>
						</Field>
						<Field orientation="horizontal">
							<Checkbox id="post-type-snapchat" />
							<FieldLabel htmlFor="post-type-snapchat">
								<SnapchatIcon className="w-4" />
								Instagram Post
							</FieldLabel>
						</Field>
						<Field orientation="horizontal">
							<Checkbox id="post-type-vsco" />
							<FieldLabel htmlFor="post-type-vsco">
								<VSCOIcon className="w-4" />
								VSCO Post
							</FieldLabel>
						</Field>
					</FieldGroup>
				</FieldSet>
			</FieldGroup>
		</form>
	);
}

function History() {
	const navigate = useNavigate({ from: Route.fullPath });
	const { username } = useUser();

	useEffect(() => {
		if (username === null) {
			navigate({ to: "/", replace: true });
		}
	}, [navigate, username]);

	return (
		<CardContent>
			<HistorySearchForm />
		</CardContent>
	);
}
