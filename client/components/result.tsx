import { timestampDate } from "@bufbuild/protobuf/wkt";
import { GalleryHorizontalIcon, Grid3x3Icon, TextAlignJustifyIcon } from "lucide-react";

import type { ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import { FileDisplay, FilesCarousel } from "@/components/file-display";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useUser } from "@/hooks/user-provider";

export function Result({ result }: { result: ScrapeResponse }) {
	const { username } = useUser();
	return username === null ? null : (
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
						{result.files.map((file) => (
							<AccordionItem
								className="flex w-full flex-col items-center"
								key={`accordion-file-${result.postType}-${result.postOwner}-${result.post}-${file}`}
								value={file}
							>
								<AccordionTrigger>
									<Label>{file}</Label>
								</AccordionTrigger>
								<AccordionContent className="max-w-[50vw]">
									<FileDisplay file={file} post={result} username={username} />
								</AccordionContent>
							</AccordionItem>
						))}
					</Accordion>
				</TabsContent>
				<TabsContent value="grid" className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
					{result.files.map((file) => (
						<FileDisplay
							key={`grid-file-${result.postType}-${result.postOwner}-${result.post}-${file}`}
							file={file}
							post={result}
							username={username}
						/>
					))}
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
