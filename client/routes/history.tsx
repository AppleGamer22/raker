import { useQuery } from "@connectrpc/connect-query";
import { useForm } from "@tanstack/react-form";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import z from "zod";

import { getUserCategories } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { PostType } from "@/buf/raker/v1/raker_pb";
import { CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { FieldGroup, FieldLegend, Field, FieldSet, FieldLabel, FieldContent, FieldTitle } from "@/components/ui/field";
import { Separator } from "@/components/ui/separator";
import { InstagramIcon } from "@/components/ui/svgs/instagram";
import { SnapchatIcon } from "@/components/ui/svgs/snapchat";
import { TikTokIcon } from "@/components/ui/svgs/tiktok";
import { VSCOIcon } from "@/components/ui/svgs/vsco";
import { Switch } from "@/components/ui/switch";
import { useUser } from "@/hooks/user-provider";

export const Route = createFileRoute("/history")({
	component: History,
});

function HistoryPostTypeForm({
	types,
	onChangeTypes,
}: {
	types: PostType[];
	onChangeTypes: (types: PostType[]) => void;
}) {
	const form = useForm({
		defaultValues: {
			types,
		},
		validators: {
			onChange: z.object({
				types: z.array(z.number()),
			}),
		},
		onSubmit: ({ value }) => onChangeTypes(value.types),
	});

	return (
		<form>
			<FieldGroup>
				<FieldSet>
					<FieldLegend>Post Types</FieldLegend>
					<form.Field name="types" mode="array">
						{(field) => (
							<FieldGroup className="flex flex-row flex-wrap gap-1 *:w-auto">
								<FieldLabel htmlFor="post-type-instagram" className="max-w-fit">
									<Field orientation="horizontal">
										<Checkbox
											id="post-type-instagram"
											onCheckedChange={(checked) => {
												if (checked) {
													field.pushValue(PostType.Instagram);
												} else {
													const index = field.state.value.indexOf(PostType.Instagram);
													if (index > -1) {
														field.removeValue(index);
													}
												}
												form.handleSubmit();
											}}
										/>
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
										<Checkbox
											id="post-type-highlight"
											onCheckedChange={(checked) => {
												if (checked) {
													field.pushValue(PostType.Highlight);
												} else {
													const index = field.state.value.indexOf(PostType.Highlight);
													if (index > -1) {
														field.removeValue(index);
													}
												}
												form.handleSubmit();
											}}
										/>
										<FieldTitle>
											<InstagramIcon className="w-4" />
											Highlight
										</FieldTitle>
									</Field>
								</FieldLabel>
								<FieldLabel htmlFor="post-type-story" className="max-w-fit">
									<Field orientation="horizontal">
										<Checkbox
											id="post-type-story"
											onCheckedChange={(checked) => {
												if (checked) {
													field.pushValue(PostType.Story);
												} else {
													const index = field.state.value.indexOf(PostType.Story);
													if (index > -1) {
														field.removeValue(index);
													}
												}
												form.handleSubmit();
											}}
										/>
										<FieldTitle>
											<InstagramIcon className="w-4" />
											Story
										</FieldTitle>
									</Field>
								</FieldLabel>
								<FieldLabel htmlFor="post-type-tiktok" className="max-w-fit">
									<Field orientation="horizontal">
										<Checkbox
											id="post-type-tiktok"
											onCheckedChange={(checked) => {
												if (checked) {
													field.pushValue(PostType.TikTok);
												} else {
													const index = field.state.value.indexOf(PostType.TikTok);
													if (index > -1) {
														field.removeValue(index);
													}
												}
												form.handleSubmit();
											}}
										/>
										<FieldTitle>
											<TikTokIcon className="w-4" />
											Post
										</FieldTitle>
									</Field>
								</FieldLabel>
								<FieldLabel htmlFor="post-type-snapchat" className="max-w-fit">
									<Field orientation="horizontal">
										<Checkbox
											id="post-type-snapchat"
											onCheckedChange={(checked) => {
												if (checked) {
													field.pushValue(PostType.Snapchat);
												} else {
													const index = field.state.value.indexOf(PostType.Snapchat);
													if (index > -1) {
														field.removeValue(index);
													}
												}
												form.handleSubmit();
											}}
										/>
										<FieldTitle>
											<SnapchatIcon className="w-4" />
											Highlight
										</FieldTitle>
									</Field>
								</FieldLabel>
								<FieldLabel htmlFor="post-type-vsco" className="max-w-fit">
									<Field orientation="horizontal">
										<Checkbox
											id="post-type-vsco"
											onCheckedChange={(checked) => {
												if (checked) {
													field.pushValue(PostType.VSCO);
												} else {
													const index = field.state.value.indexOf(PostType.VSCO);
													if (index > -1) {
														field.removeValue(index);
													}
												}
												form.handleSubmit();
											}}
										/>
										<FieldTitle>
											<VSCOIcon className="w-4" />
											Post
										</FieldTitle>
									</Field>
								</FieldLabel>
							</FieldGroup>
						)}
					</form.Field>
				</FieldSet>
			</FieldGroup>
		</form>
	);
}

function HistoryPostCategoryForm({
	exclusive,
	setExclusive,
	selectedCategories,
	availableCategories,
	setCategories,
}: {
	exclusive: boolean;
	setExclusive: (b: boolean) => void;
	selectedCategories: string[];
	availableCategories: string[];
	setCategories: (c: string[]) => void;
}) {
	const form = useForm({
		defaultValues: {
			exclusive,
			categories: selectedCategories,
		},
		validators: {
			onChange: z.object({
				exclusive: z.boolean().catch(false),
				categories: z.array(z.string()),
			}),
		},
		onSubmit: ({ value }) => {
			setExclusive(value.exclusive);
			setCategories(value.categories);
		},
	});

	return (
		<form
			onSubmit={(e) => {
				e.preventDefault();
			}}
		>
			<FieldGroup>
				<FieldSet>
					<FieldLegend>Post Categories</FieldLegend>
					<FieldGroup className="flex flex-row flex-wrap gap-1 *:w-auto">
						<form.Field name="exclusive">
							{(field) => (
								<FieldLabel htmlFor="category-exclusive" className="max-w-fit">
									<Field orientation="horizontal">
										<Switch
											id="category-exclusive"
											name={field.name}
											checked={field.state.value}
											onCheckedChange={(checked) => {
												field.handleChange(checked);
												form.handleSubmit();
											}}
										/>
										<FieldContent>
											<FieldTitle>Exclusive</FieldTitle>
										</FieldContent>
									</Field>
								</FieldLabel>
							)}
						</form.Field>
						<Separator orientation="vertical" />
						<form.Field name="categories" mode="array">
							{(field) => (
								<>
									{availableCategories.map((category) => (
										<FieldLabel
											key={`category-${category}`}
											htmlFor={`category-${category}`}
											className="max-w-fit"
										>
											<Field orientation="horizontal">
												<Checkbox
													id={`category-${category}`}
													name={field.name}
													checked={field.state.value.includes(category)}
													onCheckedChange={(checked) => {
														if (checked) {
															field.pushValue(category);
														} else {
															const index = field.state.value.indexOf(category);
															if (index > -1) {
																field.removeValue(index);
															}
														}
														form.handleSubmit();
													}}
												/>
												<FieldContent>
													<FieldTitle>{category}</FieldTitle>
												</FieldContent>
											</Field>
										</FieldLabel>
									))}
								</>
							)}
						</form.Field>
					</FieldGroup>
				</FieldSet>
			</FieldGroup>
		</form>
	);
}

function History() {
	const navigate = useNavigate({ from: Route.fullPath });
	const { username } = useUser();
	const [types, setTypes] = useState([
		PostType.Instagram,
		PostType.Highlight,
		PostType.Story,
		PostType.TikTok,
		PostType.Snapchat,
		PostType.VSCO,
	]);
	const categoriesQuery = useQuery(getUserCategories, {});
	const availableCategories = categoriesQuery.data?.categories ?? [];
	const [categories, setCategories] = useState<string[]>(availableCategories);
	const [exclusive, setExclusive] = useState(false);

	useEffect(() => {
		if (username === null) {
			navigate({ to: "/", replace: true });
		}
	}, [navigate, username]);

	return (
		<CardContent>
			<HistoryPostTypeForm types={types} onChangeTypes={setTypes} />
			<Separator className="my-2" />
			<HistoryPostCategoryForm
				exclusive={exclusive}
				setExclusive={setExclusive}
				availableCategories={availableCategories}
				selectedCategories={categories}
				setCategories={setCategories}
			/>
		</CardContent>
	);
}
