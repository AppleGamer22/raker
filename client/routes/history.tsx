import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation } from "@connectrpc/connect-query";
import { useForm } from "@tanstack/react-form";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { SearchIcon } from "lucide-react";
import { Fragment, useEffect, useState } from "react";
import { toast } from "sonner";
import { z } from "zod";

import { searchHistory, searchHistoryOwners } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { PostType, type ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { FilesCarousel } from "@/components/file-display";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import {
	Combobox,
	ComboboxChip,
	ComboboxChips,
	ComboboxChipsInput,
	ComboboxContent,
	ComboboxGroup,
	ComboboxItem,
	ComboboxLabel,
	ComboboxList,
	ComboboxValue,
	useComboboxAnchor,
} from "@/components/ui/combobox";
import { FieldGroup, FieldLegend, Field, FieldSet, FieldLabel, FieldContent, FieldTitle } from "@/components/ui/field";
import { InputGroupAddon } from "@/components/ui/input-group";
import { Label } from "@/components/ui/label";
import {
	Pagination,
	PaginationContent,
	PaginationItem,
	PaginationLink,
	PaginationNext,
	PaginationPrevious,
} from "@/components/ui/pagination";
import { Progress } from "@/components/ui/progress";
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
const defaultPostTypes = [
	PostType.Instagram,
	PostType.Highlight,
	PostType.Story,
	PostType.TikTok,
	PostType.Snapchat,
	PostType.VSCO,
];

type OwnerPostType = {
	owner: string;
	type: PostType | -1;
};

const historyFormSchema = z.object({
	types: z.array(z.enum(PostType)),
	exclusive: z.boolean(),
	categories: z.array(z.string()),
	ownerSearchTerm: z.string(),
	ownersSearchValue: z.array(
		z.object({
			owner: z.string(),
			type: z.union([z.enum(PostType), z.literal(-1)]),
		}),
	),
});

type HistoryFormValues = z.infer<typeof historyFormSchema>;

function PostTypeIconLabel({ type }: { type: PostType }) {
	switch (type) {
		case PostType.Instagram:
			return (
				<span className="inline-flex items-center gap-1 whitespace-nowrap">
					<InstagramIcon className="w-4" />
					Post
				</span>
			);
		case PostType.Highlight:
			return (
				<span className="inline-flex items-center gap-1 whitespace-nowrap">
					<InstagramIcon className="w-4" />
					Highlight
				</span>
			);
		case PostType.Story:
			return (
				<span className="inline-flex items-center gap-1 whitespace-nowrap">
					<InstagramIcon className="w-4" />
					Story
				</span>
			);
		case PostType.TikTok:
			return (
				<span className="inline-flex items-center gap-1 whitespace-nowrap">
					<TikTokIcon className="w-4" />
					Post
				</span>
			);
		case PostType.Snapchat:
			return (
				<span className="inline-flex items-center gap-1 whitespace-nowrap">
					<SnapchatIcon className="w-4" />
					Highlight
				</span>
			);
		case PostType.VSCO:
			return (
				<span className="inline-flex items-center gap-1 whitespace-nowrap">
					<VSCOIcon className="w-4" />
					Post
				</span>
			);
	}
}

function PlatformIcon({ type }: { type: PostType | -1 }) {
	switch (type) {
		case PostType.Instagram:
		case PostType.Highlight:
		case PostType.Story:
			return <InstagramIcon className="w-4" />;
		case PostType.TikTok:
			return <TikTokIcon className="w-4" />;
		case PostType.Snapchat:
			return <SnapchatIcon className="w-4" />;
		case PostType.VSCO:
			return <VSCOIcon className="w-4" />;
	}
}

function HistoryPagination({
	current,
	total,
	onChange,
}: {
	current: number;
	total: number;
	onChange: (n: number) => void;
}) {
	return total <= 1 ? null : (
		<Pagination>
			<PaginationContent>
				{current > 0 && (
					<PaginationItem>
						<PaginationPrevious onClick={() => onChange(current - 1)} />
					</PaginationItem>
				)}
				<PaginationItem>
					<PaginationLink>{current}</PaginationLink>
				</PaginationItem>
				{current < total - 1 && (
					<PaginationItem>
						<PaginationNext onClick={() => onChange(current + 1)} />
					</PaginationItem>
				)}
			</PaginationContent>
		</Pagination>
	);
}

function History() {
	const navigate = useNavigate({ from: Route.fullPath });
	const { username, categories: availableCategories } = useUser();
	const [ownersSearchOptions, setOwnersSearchOptions] = useState<OwnerPostType[]>([]);
	const [isOpen, setIsOpen] = useState(false);
	const [totalCount, setTotalCount] = useState(BigInt(0));
	const [currentPage, setCurrentPage] = useState(0);
	const [pageCount] = useState(0);
	const [histories, setHistories] = useState<ScrapeResponse[]>([]);

	const anchor = useComboboxAnchor();

	const ownersMutation = useMutation(searchHistoryOwners);
	const searchHistoryMutation = useMutation(searchHistory);
	const form = useForm({
		defaultValues: {
			types: defaultPostTypes,
			exclusive: false,
			categories: availableCategories,
			ownerSearchTerm: "",
			ownersSearchValue: [],
		} as HistoryFormValues,
		validators: {
			onChange: historyFormSchema,
			onSubmit: historyFormSchema,
		},
		onSubmit: async ({ value: { categories, exclusive, ownersSearchValue, types } }) => {
			try {
				const { histories, totalCount } = await searchHistoryMutation.mutateAsync({
					categories: categories,
					exclusive: exclusive,
					types: types,
					owners: ownersSearchValue.map(({ owner }) => owner),
					page: BigInt(1),
					pageSize: 30,
				});
				setTotalCount(totalCount);
				setHistories(histories);
			} catch (err) {
				toast.error((err as Error).message, {
					position: "top-center",
				});
			}
		},
	});

	const HistoryPostTypeForm = () => {
		const typeOptions = [
			{ id: "post-type-instagram", value: PostType.Instagram, label: "Post", Icon: InstagramIcon },
			{ id: "post-type-highlight", value: PostType.Highlight, label: "Highlight", Icon: InstagramIcon },
			{ id: "post-type-story", value: PostType.Story, label: "Story", Icon: InstagramIcon },
			{ id: "post-type-tiktok", value: PostType.TikTok, label: "Post", Icon: TikTokIcon },
			{ id: "post-type-snapchat", value: PostType.Snapchat, label: "Highlight", Icon: SnapchatIcon },
			{ id: "post-type-vsco", value: PostType.VSCO, label: "Post", Icon: VSCOIcon },
		] as const;

		return (
			<FieldGroup>
				<FieldSet>
					<FieldLegend>Post Types</FieldLegend>
					<form.Field name="types" mode="array">
						{(field) => (
							<FieldGroup className="flex flex-row flex-wrap gap-1 *:w-auto">
								{typeOptions.map(({ id, value, label, Icon }) => (
									<FieldLabel key={id} htmlFor={id} className="max-w-fit">
										<Field orientation="horizontal">
											<Checkbox
												id={id}
												checked={field.state.value.includes(value)}
												onCheckedChange={(checked) => {
													if (checked) {
														if (!field.state.value.includes(value)) {
															field.pushValue(value);
														}
													} else {
														const index = field.state.value.indexOf(value);
														if (index > -1) {
															field.removeValue(index);
														}
													}
												}}
											/>
											<FieldContent>
												<FieldTitle>
													<Icon className="w-4" />
													{label}
												</FieldTitle>
											</FieldContent>
										</Field>
									</FieldLabel>
								))}
							</FieldGroup>
						)}
					</form.Field>
				</FieldSet>
			</FieldGroup>
		);
	};

	const HistoryPostCategoryForm = () => (
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
	);

	useEffect(() => {
		form.setFieldValue("categories", availableCategories);
	}, [availableCategories, form]);

	useEffect(() => {
		if (username === null) {
			navigate({ to: "/", replace: true });
		}
	}, [navigate, username]);

	return (
		<CardContent>
			<form
				onSubmit={(e) => {
					e.preventDefault();
					e.stopPropagation();
					form.handleSubmit();
				}}
			>
				<Collapsible open={isOpen} onOpenChange={setIsOpen}>
					<form.Subscribe selector={(state) => state.values}>
						{({ types, exclusive, categories }) => (
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
									<Badge variant={exclusive ? "default" : "outline"}>
										Exclusive: {exclusive ? "On" : "Off"}
									</Badge>
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
						)}
					</form.Subscribe>
					<CollapsibleContent className="mt-1">
						<HistoryPostTypeForm />
						<Separator className="my-2" />
						<HistoryPostCategoryForm />
					</CollapsibleContent>
				</Collapsible>

				<form.Field name="ownerSearchTerm">
					{(searchField) => (
						<form.Field name="ownersSearchValue">
							{(ownersField) => (
								<Combobox
									multiple
									items={
										searchField.state.value.length > 0
											? [searchField.state.value, ...ownersSearchOptions]
											: ownersSearchOptions
									}
									value={ownersField.state.value}
									onValueChange={(value) => {
										if (value !== null) {
											ownersField.handleChange(value);
											searchField.handleChange("");
										}
									}}
								>
									<ComboboxChips className="my-2" ref={anchor}>
										<ComboboxValue>
											{(values) => (
												<Fragment>
													{values.map(({ owner, type }: OwnerPostType) => (
														<ComboboxChip
															key={`search-chip-${type}-${owner}`}
															className="select-text!"
														>
															<PlatformIcon type={type} />
															{owner}
														</ComboboxChip>
													))}
													<ComboboxChipsInput
														placeholder="post owner search"
														value={searchField.state.value}
														onChange={async (e) => {
															let ownerSearchQuery = e.target.value;
															if (
																searchField.state.value.substring(0, 4) !==
																ownerSearchQuery.substring(0, 4)
															) {
																ownerSearchQuery = ownerSearchQuery.substring(0, 4);
															}
															searchField.handleChange(e.target.value);
															if (ownerSearchQuery.length === 4) {
																try {
																	const { owners } = await ownersMutation.mutateAsync(
																		{
																			categories:
																				form.getFieldValue("categories"),
																			exclusive: form.getFieldValue("exclusive"),
																			types: form.getFieldValue("types"),
																			owner: ownerSearchQuery,
																		},
																	);
																	setOwnersSearchOptions(owners);
																} catch (err) {
																	toast.error((err as Error).message, {
																		position: "top-center",
																	});
																}
															} else if (ownerSearchQuery.length === 0) {
																setOwnersSearchOptions([]);
															}
														}}
													/>
													<InputGroupAddon>
														<SearchIcon />
													</InputGroupAddon>
												</Fragment>
											)}
										</ComboboxValue>
									</ComboboxChips>
									{searchField.state.value.length > 0 && (
										<ComboboxContent>
											<ComboboxList>
												<ComboboxGroup>
													<ComboboxLabel>Search Term</ComboboxLabel>
													<ComboboxItem
														key="search-term"
														value={
															{
																owner: searchField.state.value,
																type: -1,
															} as OwnerPostType
														}
													>
														{searchField.state.value}
													</ComboboxItem>
												</ComboboxGroup>
												{ownersSearchOptions.length > 0 && (
													<ComboboxGroup>
														<ComboboxLabel>Post Owners</ComboboxLabel>
														{ownersSearchOptions
															.filter(
																(item1) =>
																	item1.owner.includes(searchField.state.value) &&
																	ownersField.state.value.filter(
																		(item2) =>
																			item2.owner === item1.owner &&
																			item2.type === item1.type,
																	).length === 0,
															)
															.map((item) => (
																<ComboboxItem
																	key={`search-${item.type}-${item.owner}`}
																	value={item}
																>
																	<PlatformIcon type={item.type} />
																	{item.owner}
																</ComboboxItem>
															))}
													</ComboboxGroup>
												)}
											</ComboboxList>
										</ComboboxContent>
									)}
								</Combobox>
							)}
						</form.Field>
					)}
				</form.Field>

				<Button className="w-full" type="submit" disabled={searchHistoryMutation.isPending}>
					Search
				</Button>
			</form>
			{searchHistoryMutation.isPending && <Progress className="pt-2" value={null} />}
			{totalCount > 0 && <Label className="my-2 justify-center">{totalCount} results</Label>}
			<HistoryPagination current={currentPage} total={pageCount} onChange={setCurrentPage} />
			<div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
				{histories.map(({ postType, postOwner, post, postDate, categories, files }) => (
					<Card key={`post-${postType}-${postOwner}-${post}`}>
						<CardHeader className="w-full wrap-break-word">
							<span className="inline-block space-x-1 leading-none *:my-0.5 *:align-middle">
								<Badge variant="secondary">
									<PostTypeIconLabel type={postType} />
								</Badge>
								<span>/</span>
								<Badge variant="secondary">
									<code className="leading-none select-text!">{postOwner}</code>
								</Badge>
								<span>/</span>
								<Badge variant="secondary">
									<code className="leading-none select-text!">{post}</code>
								</Badge>
							</span>
							{postDate !== undefined && <p>{timestampDate(postDate).toString()}</p>}
							<span>
								{categories.map((category) => (
									<Badge
										key={`category-${postType}-${postOwner}-${post}-${category}`}
										variant="secondary"
									>
										{category}
									</Badge>
								))}
							</span>
						</CardHeader>
						<CardContent>
							<FilesCarousel
								post={{ postType, postOwner, post, files } as ScrapeResponse}
								username={username!}
							/>
						</CardContent>
					</Card>
				))}
			</div>
			<HistoryPagination current={currentPage} total={pageCount} onChange={setCurrentPage} />
		</CardContent>
	);
}
