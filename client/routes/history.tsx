import { useMutation } from "@connectrpc/connect-query";
import { useForm } from "@tanstack/react-form";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { SearchIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import z from "zod";

import { searchHistoryOwners } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { PostType } from "@/buf/raker/v1/raker_pb";
import { Badge } from "@/components/ui/badge";
import { CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import {
	Combobox,
	ComboboxContent,
	ComboboxGroup,
	ComboboxInput,
	ComboboxItem,
	ComboboxLabel,
	ComboboxList,
} from "@/components/ui/combobox";
import { FieldGroup, FieldLegend, Field, FieldSet, FieldLabel, FieldContent, FieldTitle } from "@/components/ui/field";
import { InputGroupAddon } from "@/components/ui/input-group";
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

function PostTypeIconLabel({ type }: { type: PostType }) {
	switch (type) {
		case PostType.Instagram:
			return (
				<>
					<InstagramIcon className="w-4" />
					Post
				</>
			);
		case PostType.Highlight:
			return (
				<>
					<InstagramIcon className="w-4" />
					Highlight
				</>
			);
		case PostType.Story:
			return (
				<>
					<InstagramIcon className="w-4" />
					Story
				</>
			);
		case PostType.TikTok:
			return (
				<>
					<TikTokIcon className="w-4" />
					Post
				</>
			);
		case PostType.Snapchat:
			return (
				<>
					<SnapchatIcon className="w-4" />
					Highlight
				</>
			);
		case PostType.VSCO:
			return (
				<>
					<VSCOIcon className="w-4" />
					Post
				</>
			);
	}
}

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
											checked={field.state.value.includes(PostType.Instagram)}
											onCheckedChange={(checked) => {
												if (checked) {
													if (!field.state.value.includes(PostType.Instagram)) {
														field.pushValue(PostType.Instagram);
													}
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
											checked={field.state.value.includes(PostType.Highlight)}
											onCheckedChange={(checked) => {
												if (checked) {
													if (!field.state.value.includes(PostType.Highlight)) {
														field.pushValue(PostType.Highlight);
													}
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
											checked={field.state.value.includes(PostType.Story)}
											onCheckedChange={(checked) => {
												if (checked) {
													if (!field.state.value.includes(PostType.Story)) {
														field.pushValue(PostType.Story);
													}
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
											checked={field.state.value.includes(PostType.TikTok)}
											onCheckedChange={(checked) => {
												if (checked) {
													if (!field.state.value.includes(PostType.TikTok)) {
														field.pushValue(PostType.TikTok);
													}
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
											checked={field.state.value.includes(PostType.Snapchat)}
											onCheckedChange={(checked) => {
												if (checked) {
													if (!field.state.value.includes(PostType.Snapchat)) {
														field.pushValue(PostType.Snapchat);
													}
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
											checked={field.state.value.includes(PostType.VSCO)}
											onCheckedChange={(checked) => {
												if (checked) {
													if (!field.state.value.includes(PostType.VSCO)) {
														field.pushValue(PostType.VSCO);
													}
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
															if (!field.state.value.includes(category)) {
																field.pushValue(category);
															}
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
	const { username, categories: availableCategories } = useUser();
	const [types, setTypes] = useState([
		PostType.Instagram,
		PostType.Highlight,
		PostType.Story,
		PostType.TikTok,
		PostType.Snapchat,
		PostType.VSCO,
	]);
	const [categories, setCategories] = useState<string[]>(availableCategories);
	const [ownersSearchOptions, setOwnersSearchOptions] = useState<string[]>([]);
	const [exclusive, setExclusive] = useState(false);
	const [isOpen, setIsOpen] = useState(false);
	const [ownerSearchTerm, setOwnerSearchTerm] = useState("");
	const [ownerSearchValue, setOwnerSearchValue] = useState("");
	const comboboxItems = ownerSearchTerm.length > 0 ? [ownerSearchTerm, ...ownersSearchOptions] : ownersSearchOptions;

	const ownersMutation = useMutation(searchHistoryOwners);

	useEffect(() => {
		setCategories(availableCategories);
	}, [availableCategories]);

	useEffect(() => {
		if (username === null) {
			navigate({ to: "/", replace: true });
		}
	}, [navigate, username]);

	return (
		<CardContent>
			<Collapsible open={isOpen} onOpenChange={setIsOpen}>
				<CollapsibleTrigger className="w-full rounded-md border px-3 py-2 text-left hover:bg-muted/40">
					<div className="flex flex-wrap items-center gap-2">
						{types.length > 0 ? (
							types.map((type, index) => (
								<Badge key={`type-summary-${type}-${index}`} variant="secondary">
									<PostTypeIconLabel type={type} />
								</Badge>
							))
						) : (
							<Badge variant="ghost">No post types selected</Badge>
						)}
					</div>
					<Separator className="my-2" />
					<div className="flex flex-wrap items-center gap-2">
						<Badge variant={exclusive ? "default" : "outline"}>Exclusive: {exclusive ? "On" : "Off"}</Badge>
						{categories.length > 0 ? (
							categories.map((category, index) => (
								<Badge key={`category-summary-${category}-${index}`} variant="default">
									{category}
								</Badge>
							))
						) : (
							<Badge variant="ghost">No categories selected</Badge>
						)}
					</div>
				</CollapsibleTrigger>
				<CollapsibleContent>
					<HistoryPostTypeForm types={types} onChangeTypes={setTypes} />
					<Separator className="my-2" />
					<HistoryPostCategoryForm
						exclusive={exclusive}
						setExclusive={setExclusive}
						availableCategories={availableCategories}
						selectedCategories={categories}
						setCategories={setCategories}
					/>
				</CollapsibleContent>
			</Collapsible>
			<Combobox
				items={comboboxItems}
				value={ownerSearchValue}
				onValueChange={(value) => (value !== null ? setOwnerSearchValue(value) : null)}
			>
				<ComboboxInput
					className="my-2"
					placeholder="post owner search"
					showClear
					onChange={async (e) => {
						setOwnerSearchTerm(e.target.value);
						if (e.target.value.length === 4) {
							try {
								const { owners } = await ownersMutation.mutateAsync({
									categories,
									exclusive,
									types,
									owner: e.target.value,
								});
								setOwnersSearchOptions(owners);
							} catch (err) {
								toast.error((err as Error).message, {
									position: "top-center",
								});
							}
						} else if (e.target.value.length === 0) {
							setOwnersSearchOptions([]);
						}
					}}
				>
					<InputGroupAddon>
						<SearchIcon />
					</InputGroupAddon>
				</ComboboxInput>
				{ownerSearchTerm.length > 0 && (
					<ComboboxContent>
						<ComboboxList>
							{ownerSearchTerm.length > 0 && (
								<ComboboxGroup>
									<ComboboxLabel>Search Term</ComboboxLabel>
									<ComboboxItem key="search-term" value={ownerSearchTerm}>
										{ownerSearchTerm}
									</ComboboxItem>
								</ComboboxGroup>
							)}
							{ownersSearchOptions.length > 0 && (
								<ComboboxGroup>
									<ComboboxLabel>Owners</ComboboxLabel>
									{ownersSearchOptions
										.filter((item) => item.includes(ownerSearchTerm))
										.map((item) => (
											<ComboboxItem key={`search-${item}`} value={item}>
												{item}
											</ComboboxItem>
										))}
								</ComboboxGroup>
							)}
						</ComboboxList>
					</ComboboxContent>
				)}
			</Combobox>
		</CardContent>
	);
}
