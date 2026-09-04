import { useMutation } from "@connectrpc/connect-query";
import { useForm } from "@tanstack/react-form";
import { createFileRoute, Link, stripSearchParams, useNavigate } from "@tanstack/react-router";
import { EllipsisIcon, SearchIcon } from "lucide-react";
import {
	Fragment,
	useEffect,
	useId,
	useRef,
	useState,
	type Dispatch,
	type ReactNode,
	type SetStateAction,
} from "react";
import { z } from "zod";

import { searchHistory, searchHistoryOwners } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import { PostType, type ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { FilesCarousel } from "@/components/file-display";
import {
	EditHistoryCategoriesForm,
	PlatformIcon,
	PostOwnerContextMenuContext,
	PostTypeIconLabel,
	ResultHeader,
} from "@/components/result";
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
import { ContextMenu, ContextMenuTrigger } from "@/components/ui/context-menu";
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
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { InstagramIcon } from "@/components/ui/svgs/instagram";
import { SnapchatIcon } from "@/components/ui/svgs/snapchat";
import { TikTokIcon } from "@/components/ui/svgs/tiktok";
import { VSCOIcon } from "@/components/ui/svgs/vsco";
import { Switch } from "@/components/ui/switch";
import { useUser } from "@/hooks/user-provider";
import { timestampFormat, Toaster } from "@/lib/utils";
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
	legendBadge,
	formPrefix,
	exclusiveField,
	categoriesField,
}: {
	availableCategories: string[];
	showExclusive?: boolean;
	legendBadge?: ReactNode;
	formPrefix?: string;
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
}) {
	const prefix = formPrefix ?? useId();

	return (
		<FieldGroup>
			<FieldSet>
				<FieldLegend className="flex items-center">
					{legendBadge}
					Post Categories
				</FieldLegend>
				<FieldGroup className="flex flex-row flex-wrap gap-1 *:w-auto">
					{showExclusive && exclusiveField ? (
						<>
							<FieldLabel htmlFor={`${prefix}-category-exclusive`} className="max-w-fit">
								<Field orientation="horizontal">
									<Switch
										id={`${prefix}-category-exclusive`}
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
							<FieldLabel htmlFor={`${prefix}-category-only-video`} className="max-w-fit">
								<Field orientation="horizontal">
									<Switch id={`${prefix}-category-only-video`} />
									<FieldContent>
										<FieldTitle>Only Video</FieldTitle>
									</FieldContent>
								</Field>
							</FieldLabel>
							<Separator orientation="vertical" />
						</>
					) : null}
					{availableCategories.map((category) => (
						<FieldLabel
							key={`category-${category}`}
							htmlFor={`${prefix}-category-${category}`}
							className="max-w-fit"
						>
							<Field orientation="horizontal">
								<Checkbox
									id={`${prefix}-category-${category}`}
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

function HistoryPostTypeForm({
	formPrefix,
	typesField,
}: {
	formPrefix?: string;
	typesField: {
		name: string;
		value: HistoryFormValues["types"];
		onToggleType: (type: PostType, checked: boolean) => void;
	};
}) {
	const prefix = formPrefix ?? useId();

	return (
		<FieldGroup>
			<FieldSet>
				<FieldLegend>Post Types</FieldLegend>
				<FieldGroup className="flex flex-row flex-wrap gap-1 *:w-auto">
					{postTypeOptions.map(({ id, value, label, Icon }) => (
						<FieldLabel key={id} htmlFor={`${prefix}-${id}`} className="max-w-fit">
							<Field orientation="horizontal">
								<Checkbox
									id={`${prefix}-${id}`}
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

function HistoryCard({
	history,
	username,
	exclusive,
	linkTarget,
}: {
	history: ScrapeResponse;
	username: string | undefined;
	exclusive: boolean;
	linkTarget: string | undefined;
}) {
	const { categories: availableCategories } = useUser();
	const [result, setResult] = useState(history);
	const [open, setOpen] = useState(false);

	return (
		<Card key={`post-${history.postType}-${history.postOwner}-${history.post}`}>
			<CardHeader className="w-full wrap-break-word">
				<div className="flex max-w-full flex-wrap items-center gap-x-1 gap-y-1 leading-none">
					<ResultHeader result={history} categories={availableCategories} exclusive={exclusive} showPost />
					<label className="basis-full">{timestampFormat(history.postDate!)}</label>
					<span className="inline-flex flex-wrap items-center gap-1">
						{result.categories.map((category) => (
							<Badge
								key={`category-${history.postType}-${history.postOwner}-${history.post}-${category}`}
								variant="secondary"
							>
								<Link
									to="/history"
									search={{
										categories: [category],
										exclusive: exclusive,
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
						<Sheet open={open} onOpenChange={setOpen}>
							<SheetTrigger>
								<Badge
									className="cursor-pointer hover:bg-secondary/80"
									variant="secondary"
									onClick={() => setOpen(true)}
								>
									<EllipsisIcon />
								</Badge>
							</SheetTrigger>
							<SheetContent side="bottom">
								<EditHistoryCategoriesForm
									availableCategories={availableCategories}
									result={result}
									setResult={(result) => {
										setResult(result);
										setOpen(false);
									}}
								/>
							</SheetContent>
						</Sheet>
					</span>
				</div>
			</CardHeader>
			<CardContent>
				<FilesCarousel post={result} username={username!} />
			</CardContent>
		</Card>
	);
}

function serializeSearchParams({
	types,
	exclusive,
	categories,
	owners,
	page,
}: {
	types: PostType[];
	exclusive: boolean;
	categories: string[];
	owners?: OwnerPostType[];
	page: bigint;
}) {
	return JSON.stringify({
		types: [...types].sort(),
		exclusive,
		categories: [...categories].sort(),
		owners: (owners ?? [])
			.map(({ owner, type }) => ({ owner, type }))
			.sort((a, b) => a.owner.localeCompare(b.owner) || a.type - b.type),
		page: page.toString(),
	});
}

export function HistorySearchForm({
	owners = [],
	types = defaultPostTypes,
	exclusive = false,
	categories = [],
	currentPage = 1n,
	pageSize = 30,
	autoSubmit = true,
	setHistories,
	setTotalCount,
	setHistorySearchPending,
	setExclusive,
	onSuccess,
	onError,
	onSearchSubmit,
	onResult,
}: {
	owners?: OwnerPostType[];
	types?: PostType[];
	exclusive?: boolean;
	categories?: string[];
	currentPage?: bigint;
	pageSize?: number;
	autoSubmit?: boolean;
	setHistories?: Dispatch<SetStateAction<ScrapeResponse[]>>;
	setTotalCount?: Dispatch<SetStateAction<bigint>>;
	setHistorySearchPending?: Dispatch<SetStateAction<boolean>>;
	setExclusive?: Dispatch<SetStateAction<boolean>>;
	onSuccess?: (result: { histories: ScrapeResponse[]; totalCount: bigint }) => void;
	onError?: (error: Error) => void;
	onSearchSubmit?: (values: {
		categories: string[];
		exclusive: boolean;
		ownersSearchValue: OwnerPostType[];
		types: PostType[];
		page: bigint;
	}) => void;
	onResult?: (formValue: {
		categories: string[];
		exclusive: boolean;
		ownersSearchValue: OwnerPostType[];
		types: PostType[];
	}) => void;
}) {
	const { username, isCategoriesPending, categories: availableCategories } = useUser();

	const ownersMutation = useMutation(searchHistoryOwners);
	const searchHistoryMutation = useMutation(searchHistory);

	const [isOpen, setIsOpen] = useState(false);
	const [ownersSearchOptions, setOwnersSearchOptions] = useState<OwnerPostType[]>(owners ?? []);

	const anchor = useComboboxAnchor();
	const lastExecutedSearch = useRef<string>("");

	const executeSearch = async (searchParams: {
		categories: string[];
		exclusive: boolean;
		types: PostType[];
		owners: OwnerPostType[];
		page: bigint;
	}) => {
		const key = serializeSearchParams({
			types: searchParams.types,
			exclusive: searchParams.exclusive,
			categories: searchParams.categories,
			owners: searchParams.owners,
			page: searchParams.page,
		});

		lastExecutedSearch.current = key;
		setHistorySearchPending?.(true);

		const handleSearchSubmit = onSearchSubmit ?? onResult;
		handleSearchSubmit?.({
			categories: searchParams.categories,
			exclusive: searchParams.exclusive,
			ownersSearchValue: searchParams.owners,
			types: searchParams.types,
			page: searchParams.page,
		});

		try {
			const { histories, totalCount } = await searchHistoryMutation.mutateAsync({
				categories: searchParams.categories,
				exclusive: searchParams.exclusive,
				types: searchParams.types,
				owners: searchParams.owners.map(({ owner }) => owner),
				page: searchParams.page,
				pageSize,
			});

			setTotalCount?.(totalCount);
			setHistories?.(histories);
			setExclusive?.(searchParams.exclusive);
			onSuccess?.({ histories, totalCount });
		} catch (err) {
			onError?.(err as Error);
			Toaster.error(err as Error);
		} finally {
			setHistorySearchPending?.(false);
		}
	};

	const form = useForm({
		defaultValues: {
			types,
			exclusive,
			categories,
			ownerSearchTerm: "",
			ownersSearchValue: owners ?? [],
		} as HistoryFormValues,
		validators: {
			onChange: historyFormSchema,
			onSubmit: historyFormSchema,
		},
		onSubmit: async ({ value: { categories, exclusive, ownersSearchValue, types } }) => {
			await executeSearch({
				categories,
				exclusive,
				types,
				owners: ownersSearchValue,
				page: currentPage,
			});
		},
	});

	const validSearchCategories = categories.filter((category) => availableCategories.includes(category));
	const incomingSearchKey = serializeSearchParams({
		types,
		exclusive,
		categories: validSearchCategories,
		owners,
		page: currentPage,
	});

	useEffect(() => {
		if (!autoSubmit || isCategoriesPending || username === null) {
			return;
		}

		if (lastExecutedSearch.current === incomingSearchKey) {
			return;
		}

		form.setFieldValue("types", types);
		form.setFieldValue("exclusive", exclusive);
		form.setFieldValue("categories", validSearchCategories);
		form.setFieldValue("ownersSearchValue", owners ?? []);

		executeSearch({
			categories: validSearchCategories,
			exclusive,
			types,
			owners: owners ?? [],
			page: currentPage,
		});
	}, [autoSubmit, incomingSearchKey, isCategoriesPending, username]);

	return (
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
										formPrefix="search"
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
											{(values: OwnerPostType[]) => (
												<Fragment>
													{values.map(({ owner, type }) => (
														<ContextMenu key={`${owner}-${type}`}>
															<ContextMenuTrigger>
																<ComboboxChip key={`search-chip-${type}-${owner}`}>
																	<PlatformIcon type={type} />
																	{owner}
																</ComboboxChip>
															</ContextMenuTrigger>
															<PostOwnerContextMenuContext
																result={
																	{
																		postOwner: owner,
																		postType: type,
																	} as ScrapeResponse
																}
																categories={availableCategories}
																exclusive={exclusive}
															/>
														</ContextMenu>
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
																	Toaster.error(err as Error);
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
	);
}

function History() {
	const { types, exclusive, categories, owners, page } = Route.useSearch();
	const navigate = useNavigate({ from: Route.fullPath });
	const { username } = useUser();
	const linkTarget = inPWA() ? undefined : "_blank";
	const [totalCount, setTotalCount] = useState(0n);
	const [histories, setHistories] = useState<ScrapeResponse[]>([]);
	const [isSearching, setIsSearching] = useState(false);

	useEffect(() => {
		if (username === null) {
			navigate({ to: "/", replace: true });
		}
	}, [navigate, username]);

	const historyPageinationButtons = (
		<>
			{isSearching && <Progress className="pt-2" value={null} />}
			<HistoryPagination
				current={page}
				total={totalCount / 30n + (totalCount % 30n ? 1n : 0n)}
				onChange={(nextPage) => {
					navigate({
						search: (prev) => ({
							...prev,
							page: nextPage,
						}),
						replace: true,
					});
				}}
			/>
		</>
	);

	return (
		<CardContent>
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
			{totalCount > 0 && <Label className="my-2 justify-center">{totalCount} results</Label>}
			{historyPageinationButtons}
			<div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
				{histories.map((history) => (
					<HistoryCard
						key={`${history.postType}-${history.postOwner}-${history.post}`}
						history={history}
						exclusive={exclusive}
						linkTarget={linkTarget}
						username={username ?? undefined}
					/>
				))}
			</div>
			{histories.length > 0 && historyPageinationButtons}
		</CardContent>
	);
}
