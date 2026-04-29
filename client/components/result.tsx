import { timestampDate } from "@bufbuild/protobuf/wkt";
import { GalleryHorizontalIcon, Grid3x3Icon, TextAlignJustifyIcon } from "lucide-react";
import { useEffect, useState } from "react";

import type { ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { FileDisplay, FilesCarousel } from "@/components/file-display";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useUser } from "@/hooks/user-provider";
import { cn } from "@/lib/utils";

export function Result({ result }: { result: ScrapeResponse }) {
	const { username } = useUser();
	const [selection, setSelection] = useState<{ selectedFiles: string[]; anchorFile: string | null }>({
		selectedFiles: [],
		anchorFile: null,
	});
	const files = result.files;

	useEffect(() => {
		setSelection((current) => {
			const selectedFiles = files.filter((file) => current.selectedFiles.includes(file));
			const anchorFile =
				current.anchorFile !== null && files.includes(current.anchorFile)
					? current.anchorFile
					: (selectedFiles[0] ?? null);

			if (selectedFiles.length === current.selectedFiles.length && anchorFile === current.anchorFile) {
				return current;
			}

			return { selectedFiles, anchorFile };
		});
	}, [files]);

	useEffect(() => {
		const onKeyDown = (event: KeyboardEvent) => {
			if (event.key === "Escape") {
				setSelection({ selectedFiles: [], anchorFile: null });
			}
		};

		window.addEventListener("keydown", onKeyDown);
		return () => window.removeEventListener("keydown", onKeyDown);
	}, []);

	if (username === null) {
		return null;
	}

	const isSelected = (file: string) => selection.selectedFiles.includes(file);

	const toggleSelection = (file: string) => {
		setSelection((current) => {
			const isCurrentlySelected = current.selectedFiles.includes(file);
			const selectedFiles = isCurrentlySelected
				? current.selectedFiles.filter((selectedFile) => selectedFile !== file)
				: [...current.selectedFiles, file];
			const anchorFile = isCurrentlySelected
				? current.anchorFile === file
					? (selectedFiles.at(-1) ?? null)
					: current.anchorFile
				: file;

			return { selectedFiles, anchorFile };
		});
	};

	const selectRange = (file: string) => {
		setSelection((current) => {
			if (current.anchorFile === null) {
				return { selectedFiles: [file], anchorFile: file };
			}

			const anchorIndex = files.indexOf(current.anchorFile);
			const fileIndex = files.indexOf(file);

			if (anchorIndex === -1 || fileIndex === -1) {
				return { selectedFiles: [file], anchorFile: file };
			}

			const start = Math.min(anchorIndex, fileIndex);
			const end = Math.max(anchorIndex, fileIndex);

			return { selectedFiles: files.slice(start, end + 1), anchorFile: current.anchorFile };
		});
	};

	const handleSelection = (file: string, event: { shiftKey: boolean; ctrlKey: boolean; metaKey: boolean }) => {
		if (event.shiftKey) {
			selectRange(file);
			return;
		}

		if (event.ctrlKey || event.metaKey) {
			toggleSelection(file);
			return;
		}

		toggleSelection(file);
	};

	return (
		<section className="my-3 flex w-full flex-col items-center gap-3">
			<Label>{timestampDate(result.postDate!).toString()}</Label>
			<Tabs className="w-full">
				<TabsList className="mx-auto w-full sm:w-1/2">
					<TabsTrigger value="list">
						<TextAlignJustifyIcon />
						List
					</TabsTrigger>
					<TabsTrigger value="grid">
						<Grid3x3Icon />
						Grid
					</TabsTrigger>
					<TabsTrigger value="carousel">
						<GalleryHorizontalIcon />
						Carousel
					</TabsTrigger>
				</TabsList>
				<TabsContent value="list">
					<Accordion multiple>
						{files.map((file) => {
							const selected = isSelected(file);

							return (
								<AccordionItem
									aria-pressed={selected}
									className={cn(
										"flex w-full flex-col rounded-lg border border-transparent px-2 py-1 transition",
										selected && "border-primary/60 bg-primary/10 ring-2 ring-primary/30",
									)}
									key={`accordion-file-${result.postType}-${result.postOwner}-${result.post}-${file}`}
									value={file}
								>
									<div className="flex w-full items-center gap-2">
										<Checkbox
											checked={selected}
											aria-label={selected ? `Deselect ${file}` : `Select ${file}`}
											onClick={(event) => {
												event.preventDefault();
												event.stopPropagation();
												handleSelection(file, event);
											}}
										/>
										<AccordionTrigger className="flex-1 gap-2">
											<Label>{file}</Label>
										</AccordionTrigger>
									</div>
									<AccordionContent className="max-w-[50vw]">
										<FileDisplay file={file} post={result} username={username} />
									</AccordionContent>
								</AccordionItem>
							);
						})}
					</Accordion>
				</TabsContent>
				<TabsContent value="grid" className="grid grid-cols-2 gap-3 lg:grid-cols-3">
					{files.map((file) => {
						const selected = isSelected(file);

						return (
							<div
								aria-pressed={selected}
								className={cn(
									"relative rounded-xl border border-transparent p-1 transition",
									selected && "border-primary/60 bg-primary/10 ring-2 ring-primary/30",
								)}
								key={`grid-file-${result.postType}-${result.postOwner}-${result.post}-${file}`}
							>
								<Checkbox
									checked={selected}
									aria-label={selected ? `Deselect ${file}` : `Select ${file}`}
									className="absolute top-2 right-2 z-10"
									onClick={(event) => {
										event.preventDefault();
										event.stopPropagation();
										handleSelection(file, event);
									}}
								/>
								<FileDisplay file={file} post={result} username={username} />
							</div>
						);
					})}
				</TabsContent>
				<TabsContent
					value="carousel"
					className="mt-2 w-full [&_img]:max-h-[50vh] [&_img]:w-auto [&_video]:max-h-[50vh] [&_video]:w-auto"
				>
					<FilesCarousel post={result} username={username} />
				</TabsContent>
			</Tabs>
		</section>
	);
}
