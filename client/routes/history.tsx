import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation } from "@connectrpc/connect-query";
import { useForm } from "@tanstack/react-form";
import { createFileRoute, Link, stripSearchParams, useNavigate } from "@tanstack/react-router";
import { SearchIcon } from "lucide-react";
import { Fragment, useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import { z } from "zod";

import { searchHistory, searchHistoryOwners } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { PostType, type ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { FilesCarousel } from "@/components/file-display";
import { PlatformIcon, PostTypeIconLabel, ResultHeader } from "@/components/result";
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
import { Input } from "@/components/ui/input";
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
import { defaultPostTypes, inPWA } from "@/lib/utils";

const historySearchDefaults = {
	exclusive: false,
	categories: [],
	page: 1n,
	owners: [],
	types: defaultPostTypes,
};

export const Route = createFileRoute("/history")({
	component: History,
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

type HistoryPostCategoryFormProps = {
	availableCategories: string[];
	showExclusive?: boolean;
	exclusiveField?: {
		name: string;
		value: HistoryFormValues["exclusive"];
		onChange: (checked: boolean) => void;
	};
	categoriesField: {
		name: string;
		value: HistoryFormValues["categories"];
		onToggleCategory: (category: string, checked: boolean) => void;
	};
};

type HistoryPostTypeFormProps = {
	typesField: {
		name: string;
		value: HistoryFormValues["types"];
		onToggleType: (type: PostType, checked: boolean) => void;
	};
};

const postTypeOptions = [
	{ id: "post-type-instagram", value: PostType.Instagram, label: "Post", Icon: InstagramIcon },
	{ id: "post-type-highlight", value: PostType.Highlight, label: "Highlight", Icon: InstagramIcon },
	{ id: "post-type-story", value: PostType.Story, label: "Story", Icon: InstagramIcon },
	{ id: "post-type-tiktok", value: PostType.TikTok, label: "Post", Icon: TikTokIcon },
	{ id: "post-type-snapchat", value: PostType.Snapchat, label: "Highlight", Icon: SnapchatIcon },
	{ id: "post-type-vsco", value: PostType.VSCO, label: "Post", Icon: VSCOIcon },
] as const;

function HistoryPagination({
	current,
	total,
	onChange,
}: {
	current: bigint;
	total: bigint;
	onChange: (n: bigint) => void;
}) {
	const [pageValue, setPageValue] = useState(current.toString());

	useEffect(() => {
		setPageValue(current.toString());
	}, [current]);

	const commitPageValue = () => {
		if (!/^\d+$/.test(pageValue)) {
			setPageValue(current.toString());
			return;
		}

		const nextPage = BigInt(pageValue);
		if (nextPage < 1n || nextPage > total) {
			setPageValue(current.toString());
			return;
		}

		if (nextPage !== current) {
			onChange(nextPage);
		}
	};

	return total <= 1n ? null : (
		<Pagination className="my-2">
			<PaginationContent>
				{current > 1n && (
					<>
						{current > 2n && (
							<PaginationItem>
								<PaginationLink onClick={() => onChange(1n)}>1</PaginationLink>
							</PaginationItem>
						)}
						<PaginationItem>
							<PaginationPrevious onClick={() => onChange(current - 1n)} />
						</PaginationItem>
					</>
				)}
				<PaginationItem>
					<Input
						aria-label="Current page"
						className="h-8 px-1 text-center"
						onBlur={commitPageValue}
						onChange={(event) => setPageValue(event.target.value)}
						onKeyDown={(event) => {
							if (event.key === "Enter") {
								event.preventDefault();
								commitPageValue();
							}
						}}
						min="1"
						max={total.toString()}
						step="1"
						type="number"
						value={pageValue}
					/>
				</PaginationItem>
				{current < total && (
					<>
						<PaginationItem>
							<PaginationNext onClick={() => onChange(current + 1n)} />
						</PaginationItem>
						{current < total - 1n && (
							<PaginationItem>
								<PaginationLink onClick={() => onChange(total)}>{total}</PaginationLink>
							</PaginationItem>
						)}
					</>
				)}
			</PaginationContent>
		</Pagination>
	);
}

export function HistoryPostCategoryForm({
	availableCategories,
	showExclusive = true,
	exclusiveField,
	categoriesField,
}: HistoryPostCategoryFormProps) {
	return (
		<FieldGroup>
			<FieldSet>
				<FieldLegend className="flex items-center">
					{!showExclusive && <Badge className="mr-1 h-2 w-2 rounded-full p-0" />}
					Post Categories
				</FieldLegend>
				<FieldGroup className="flex flex-row flex-wrap gap-1 *:w-auto">
					{showExclusive && exclusiveField ? (
						<>
							<FieldLabel htmlFor="category-exclusive" className="max-w-fit">
								<Field orientation="horizontal">
									<Switch
										id="category-exclusive"
										name={exclusiveField.name}
										checked={exclusiveField.value}
										onCheckedChange={(checked) => {
											exclusiveField.onChange(checked);
										}}
									/>
									<FieldContent>
										<FieldTitle>Exclusive</FieldTitle>
									</FieldContent>
								</Field>
							</FieldLabel>
							<Separator orientation="vertical" />
						</>
					) : null}
					{availableCategories.map((category) => (
						<FieldLabel key={`category-${category}`} htmlFor={`category-${category}`} className="max-w-fit">
							<Field orientation="horizontal">
								<Checkbox
									id={`category-${category}`}
									name={categoriesField.name}
									checked={categoriesField.value.includes(category)}
									onCheckedChange={(checked) => {
										categoriesField.onToggleCategory(category, !!checked);
									}}
								/>
								<FieldContent>
									<FieldTitle>{category}</FieldTitle>
								</FieldContent>
							</Field>
						</FieldLabel>
					))}
				</FieldGroup>
			</FieldSet>
		</FieldGroup>
	);
}

function HistoryPostTypeForm({ typesField }: HistoryPostTypeFormProps) {
	return (
		<FieldGroup>
			<FieldSet>
				<FieldLegend>Post Types</FieldLegend>
				<FieldGroup className="flex flex-row flex-wrap gap-1 *:w-auto">
					{postTypeOptions.map(({ id, value, label, Icon }) => (
						<FieldLabel key={id} htmlFor={id} className="max-w-fit">
							<Field orientation="horizontal">
								<Checkbox
									id={id}
									name={typesField.name}
									checked={typesField.value.includes(value)}
									onCheckedChange={(checked) => {
										typesField.onToggleType(value, !!checked);
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
			</FieldSet>
		</FieldGroup>
	);
}

function History() {
	const { types, exclusive, categories, owners, page } = Route.useSearch();
	const navigate = useNavigate({ from: Route.fullPath });
	const { username, categories: availableCategories, isCategoriesPending } = useUser();
	const linkTarget = inPWA() ? undefined : "_blank";
	const [ownersSearchOptions, setOwnersSearchOptions] = useState<OwnerPostType[]>([]);
	const [isOpen, setIsOpen] = useState(false);
	const [totalCount, setTotalCount] = useState(0n);
	const [currentPage, setCurrentPage] = useState(BigInt(page));
	const currentPageRef = useRef(BigInt(page));
	const [histories, setHistories] = useState<ScrapeResponse[]>([]);
	const lastSubmittedSearch = useRef("");

	const anchor = useComboboxAnchor();

	const ownersMutation = useMutation(searchHistoryOwners);
	const searchHistoryMutation = useMutation(searchHistory);
	const form = useForm({
		defaultValues: {
			types,
			exclusive,
			categories,
			ownerSearchTerm: "",
			ownersSearchValue: owners,
		} as HistoryFormValues,
		validators: {
			onChange: historyFormSchema,
			onSubmit: historyFormSchema,
		},
		onSubmit: async ({ value: { categories, exclusive, ownersSearchValue, types } }) => {
			try {
				const { histories, totalCount } = await searchHistoryMutation.mutateAsync({
					categories,
					exclusive,
					types,
					owners: ownersSearchValue.map(({ owner }) => owner),
					page: currentPageRef.current,
					pageSize: 30,
				});
				setTotalCount(totalCount);
				setHistories(histories);
				await navigate({
					search: {
						types,
						exclusive,
						categories,
						page: currentPageRef.current,
						owners: ownersSearchValue,
					},
					replace: true,
				});
			} catch (err) {
				toast.error((err as Error).message, {
					position: "top-center",
				});
			}
		},
	});

	useEffect(() => {
		if (isCategoriesPending) {
			return;
		}

		const validSearchCategories = categories.filter((category) => availableCategories.includes(category));
		const normalizedPage = BigInt(page);
		const nextSearch = JSON.stringify({
			types,
			exclusive,
			categories: validSearchCategories,
			owners,
			page: normalizedPage.toString(),
		});

		form.setFieldValue("types", types);
		form.setFieldValue("exclusive", exclusive);
		form.setFieldValue("categories", validSearchCategories);
		form.setFieldValue("ownersSearchValue", owners);
		setCurrentPage(normalizedPage);
		currentPageRef.current = normalizedPage;

		if (username === null || lastSubmittedSearch.current === nextSearch) {
			return;
		}

		lastSubmittedSearch.current = nextSearch;
		form.handleSubmit();
	}, [availableCategories, categories, exclusive, form, isCategoriesPending, owners, page, types, username]);

	useEffect(() => {
		const normalizedPage = BigInt(page);
		setCurrentPage(normalizedPage);
		currentPageRef.current = normalizedPage;
	}, [page]);

	useEffect(() => {
		if (username === null) {
			navigate({ to: "/", replace: true });
		}
	}, [navigate, username]);

	const HistoryPageinationButtons = () => (
		<>
			{searchHistoryMutation.isPending && <Progress className="pt-2" value={null} />}
			<HistoryPagination
				current={currentPage}
				total={totalCount / 30n + (totalCount % 30n ? 1n : 0n)}
				onChange={(current) => {
					setCurrentPage(current);
					currentPageRef.current = current;
					form.handleSubmit();
				}}
			/>
		</>
	);

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
						<form.Field name="types" mode="array">
							{(typesField) => (
								<HistoryPostTypeForm
									typesField={{
										name: typesField.name,
										value: typesField.state.value,
										onToggleType: (type, checked) => {
											if (checked) {
												if (!typesField.state.value.includes(type)) {
													typesField.pushValue(type);
												}
											} else {
												const index = typesField.state.value.indexOf(type);
												if (index > -1) {
													typesField.removeValue(index);
												}
											}
										},
									}}
								/>
							)}
						</form.Field>
						<Separator className="my-2" />
						<form.Field name="exclusive">
							{(exclusiveField) => (
								<form.Field name="categories" mode="array">
									{(categoriesField) => (
										<HistoryPostCategoryForm
											availableCategories={availableCategories}
											exclusiveField={{
												name: exclusiveField.name,
												value: exclusiveField.state.value,
												onChange: exclusiveField.handleChange,
											}}
											categoriesField={{
												name: categoriesField.name,
												value: categoriesField.state.value,
												onToggleCategory: (category, checked) => {
													if (checked) {
														if (!categoriesField.state.value.includes(category)) {
															categoriesField.pushValue(category);
														}
													} else {
														const index = categoriesField.state.value.indexOf(category);
														if (index > -1) {
															categoriesField.removeValue(index);
														}
													}
												},
											}}
										/>
									)}
								</form.Field>
							)}
						</form.Field>
					</CollapsibleContent>
				</Collapsible>

				<form.Field name="ownerSearchTerm">
					{(searchField) => (
						<form.Field name="ownersSearchValue">
							{(ownersField) => {
								const searchOptions = ownersSearchOptions.filter(
									(item1) =>
										item1.owner.includes(searchField.state.value) &&
										ownersField.state.value.filter(
											(item2) => item2.owner === item1.owner && item2.type === item1.type,
										).length === 0,
								);
								const showTypedSearchQuery =
									ownersField.state.value.filter(
										({ owner, type }) => owner === searchField.state.value && type === -1,
									).length === 0;
								return (
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
																		const { owners } =
																			await ownersMutation.mutateAsync({
																				categories:
																					form.getFieldValue("categories"),
																				exclusive:
																					form.getFieldValue("exclusive"),
																				types: form.getFieldValue("types"),
																				owner: ownerSearchQuery,
																			});
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
										{searchField.state.value.length > 0 &&
										(showTypedSearchQuery || searchOptions.length > 0) ? (
											<ComboboxContent>
												<ComboboxList>
													{showTypedSearchQuery && (
														<ComboboxGroup>
															<ComboboxLabel>Search Term</ComboboxLabel>
															{
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
															}
														</ComboboxGroup>
													)}
													{searchOptions.length > 0 && (
														<ComboboxGroup>
															<ComboboxLabel>Post Owners</ComboboxLabel>
															{searchOptions.map((item) => (
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
										) : null}
									</Combobox>
								);
							}}
						</form.Field>
					)}
				</form.Field>

				<Button className="w-full" type="submit" disabled={searchHistoryMutation.isPending}>
					Search
				</Button>
			</form>
			{totalCount > 0 && <Label className="my-2 justify-center">{totalCount} results</Label>}
			<HistoryPageinationButtons />
			<div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
				{histories.map(({ postType, postOwner, post, postDate, categories, files, coordinates }) => (
					<Card key={`post-${postType}-${postOwner}-${post}`}>
						<CardHeader className="w-full wrap-break-word">
							<div className="flex max-w-full flex-wrap items-center gap-x-1 gap-y-1 leading-none">
								<ResultHeader
									result={{ postType, postOwner, post, incognito: false } as ScrapeResponse}
									categories={form.getFieldValue("categories")}
									exclusive={form.getFieldValue("exclusive")}
									showPost
								/>
								{postDate !== undefined && <p>{timestampDate(postDate).toString()}</p>}
								<span className="inline-flex flex-wrap items-center gap-1">
									{categories.map((category) => (
										<Badge
											key={`category-${postType}-${postOwner}-${post}-${category}`}
											variant="secondary"
										>
											<Link
												to="/history"
												search={{
													categories: [category],
													exclusive: form.getFieldValue("exclusive"),
													page: 1n,
													owners: [],
													types: defaultPostTypes,
												}}
												target={linkTarget}
											>
												{category}
											</Link>
										</Badge>
									))}
								</span>
							</div>
						</CardHeader>
						<CardContent>
							<FilesCarousel
								post={{ postType, postOwner, post, files, coordinates } as ScrapeResponse}
								username={username!}
							/>
						</CardContent>
					</Card>
				))}
			</div>
			{histories.length > 0 && <HistoryPageinationButtons />}
		</CardContent>
	);
}
