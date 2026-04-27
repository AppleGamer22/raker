import { useCallback, useEffect, useState } from "react";

import { PostType, type ScrapeResponse } from "@/buf/raker/v1/raker_pb";
import {
	Carousel,
	CarouselContent,
	CarouselItem,
	CarouselNext,
	CarouselPrevious,
	type CarouselApi,
} from "@/components/ui/carousel";

function postTypeString(type: PostType) {
	switch (type) {
		case PostType.Instagram:
			return "instagram";
		case PostType.Highlight:
			return "highlight";
		case PostType.Story:
			return "story";
		case PostType.TikTok:
			return "tiktok";
		case PostType.Snapchat:
			return "snapchat";
		case PostType.VSCO:
			return "vsco";
	}
}

export function FilesCarousel({
	username,
	post: { postType, postOwner, post, files },
}: {
	username: string;
	post: ScrapeResponse;
}) {
	const [api, setApi] = useState<CarouselApi>();
	const [current, setCurrent] = useState(1);

	const syncSelectedSlideHeight = useCallback((emblaApi: CarouselApi) => {
		if (!emblaApi) {
			return;
		}

		const selectedIndex = emblaApi.selectedScrollSnap();
		const selectedSlide = emblaApi.slideNodes()[selectedIndex];
		if (!selectedSlide) {
			return;
		}

		emblaApi.containerNode().style.height = `${selectedSlide.offsetHeight}px`;
	}, []);

	useEffect(() => {
		if (!api) {
			return;
		}

		const onSelect = () => {
			setCurrent(api.selectedScrollSnap() + 1);
			syncSelectedSlideHeight(api);
		};

		const onReInit = () => syncSelectedSlideHeight(api);

		// onSelect();
		api.on("reInit", onReInit);
		api.on("select", onSelect);
		api.on("settle", onSelect);

		return () => {
			api.off("reInit", onReInit);
			api.off("select", onSelect);
			api.off("settle", onSelect);
		};
	}, [api, syncSelectedSlideHeight]);

	return files.length > 1 ? (
		<Carousel opts={{ loop: true }} setApi={setApi}>
			<CarouselContent className="items-start">
				{files.map((file, i) => (
					<CarouselItem key={`file-${postType}-${postOwner}-${post}-${i}`} className="self-start">
						<FileDisplay
							username={username}
							post={{ postType, postOwner } as ScrapeResponse}
							file={file}
							onMediaLoad={() => {
								api?.reInit();
								if (api) {
									requestAnimationFrame(() => syncSelectedSlideHeight(api));
								}
							}}
						/>
					</CarouselItem>
				))}
			</CarouselContent>
			<CarouselPrevious className="top-1/2 left-2 z-10 -translate-y-1/2" />
			<CarouselNext className="top-1/2 right-2 z-10 -translate-y-1/2" />
			<div className="py-2 text-center text-sm text-muted-foreground">
				{current} of {files.length}
			</div>
		</Carousel>
	) : (
		<FileDisplay username={username} post={{ postType, postOwner } as ScrapeResponse} file={files[0]} />
	);
}

export function FileDisplay({
	username,
	file,
	post: { postType, postOwner },
	onMediaLoad,
}: {
	username: string;
	file: string;
	post: ScrapeResponse;
	onMediaLoad?: () => void;
}) {
	const url = `/api/storage/${username}/${postTypeString(postType)}/${postOwner}/${file}`;
	if (/\.(jpg)|(jpeg)|(webp)|(heic)$/.test(file)) {
		return <img src={url} onLoad={onMediaLoad} className="h-auto w-full" />;
	} else if (/\.(mp4)|(webm)$/.test(file)) {
		return <video src={url} loop controls muted onLoadedMetadata={onMediaLoad} className="h-auto w-full" />;
	} else {
		return <a href={url}>{url}</a>;
	}
}
