import { CheckboxGroup } from "@base-ui/react";
import { useQuery } from "@connectrpc/connect-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect } from "react";

import { getUserCategories } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import {
	FieldGroup,
	FieldLegend,
	Field,
	FieldSet,
	FieldLabel,
	FieldContent,
	FieldTitle,
	FieldSeparator,
} from "@/components/ui/field";
import { InstagramIcon } from "@/components/ui/svgs/instagram";
import { SnapchatIcon } from "@/components/ui/svgs/snapchat";
import { TikTokIcon } from "@/components/ui/svgs/tiktok";
import { VSCOIcon } from "@/components/ui/svgs/vsco";
import { useUser } from "@/hooks/user-provider";

export const Route = createFileRoute("/history")({
	component: History,
});

function HistoryPostTypeForm() {
	// const categoriesQuery = useQuery(getUserCategories, {});
	// const categories = categoriesQuery.data?.categories ?? [];

	return (
		<form>
			<FieldGroup>
				<FieldSet>
					<FieldLegend>Post Types</FieldLegend>
					<FieldGroup data-slot="checkbox-group">
						<CheckboxGroup className="flex flex-row flex-wrap gap-1 *:w-auto">
							<FieldLabel htmlFor="post-type-instagram" className="max-w-fit">
								<Field orientation="horizontal">
									<Checkbox id="post-type-instagram" />
									<FieldContent>
										<FieldTitle>
											<InstagramIcon className="w-4" />
											Post
										</FieldTitle>
									</FieldContent>
								</Field>
							</FieldLabel>
							<FieldLabel htmlFor="post-type-highlight" className="max-w-fit">
								<Field orientation="horizontal">
									<Checkbox id="post-type-highlight" />
									<FieldTitle>
										<InstagramIcon className="w-4" />
										Highlight
									</FieldTitle>
								</Field>
							</FieldLabel>
							<FieldLabel htmlFor="post-type-story" className="max-w-fit">
								<Field orientation="horizontal">
									<Checkbox id="post-type-story" />
									<FieldTitle>
										<InstagramIcon className="w-4" />
										Story
									</FieldTitle>
								</Field>
							</FieldLabel>
							<FieldLabel htmlFor="post-type-tiktok" className="max-w-fit">
								<Field orientation="horizontal">
									<Checkbox id="post-type-tiktok" />
									<FieldTitle>
										<TikTokIcon className="w-4" />
										Post
									</FieldTitle>
								</Field>
							</FieldLabel>
							<FieldLabel htmlFor="post-type-snapchat" className="max-w-fit">
								<Field orientation="horizontal">
									<Checkbox id="post-type-snapchat" />
									<FieldTitle>
										<SnapchatIcon className="w-4" />
										Highlight
									</FieldTitle>
								</Field>
							</FieldLabel>
							<FieldLabel htmlFor="post-type-vsco" className="max-w-fit">
								<Field orientation="horizontal">
									<Checkbox id="post-type-vsco" />
									<FieldTitle>
										<VSCOIcon className="w-4" />
										Post
									</FieldTitle>
								</Field>
							</FieldLabel>
						</CheckboxGroup>
					</FieldGroup>
				</FieldSet>
			</FieldGroup>
		</form>
	);
}

function HistoryPostCategoryForm() {
	const categoriesQuery = useQuery(getUserCategories, {});
	const categories = categoriesQuery.data?.categories ?? [];

	return (
		<form>
			<FieldGroup>
				<FieldSet>
					<FieldLegend>Post Types</FieldLegend>
					<FieldGroup data-slot="checkbox-group">
						<CheckboxGroup className="flex flex-row flex-wrap gap-1 *:w-auto">
							{categories.map((category) => (
								<FieldLabel
									key={`category-${category}`}
									htmlFor={`category-${category}`}
									className="max-w-fit"
								>
									<Field orientation="horizontal">
										<Checkbox id={`category-${category}`} />
										<FieldContent>
											<FieldTitle>{category}</FieldTitle>
										</FieldContent>
									</Field>
								</FieldLabel>
							))}
						</CheckboxGroup>
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
			<HistoryPostTypeForm />
			<FieldSeparator />
			<HistoryPostCategoryForm />
		</CardContent>
	);
}
