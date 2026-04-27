import { useEffect, useState } from "react";

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

export function PostCarousel({
	username,
	post: { postType, postOwner, post, files },
}: {
	username: string;
	post: ScrapeResponse;
}) {
	const [api, setApi] = useState<CarouselApi>();
	const [current, setCurrent] = useState(0);

	useEffect(() => {
		if (!api) {
			return;
		}

		setCurrent(api.selectedScrollSnap() + 1);

		api.on("select", () => {
			setCurrent(api.selectedScrollSnap() + 1);
		});
	}, [api]);

	return files.length > 1 ? (
		<Carousel opts={{ loop: true }} setApi={setApi} className="w-full">
			<CarouselContent>
				{files.map((file, i) => (
					<CarouselItem key={`file-${postType}-${postOwner}-${post}-${i}`}>
						<FileDisplay username={username} post={{ postType, postOwner } as ScrapeResponse} file={file} />
					</CarouselItem>
				))}
			</CarouselContent>
			<CarouselPrevious className="top-1/2 left-2 z-10 -translate-y-1/2" />
			<CarouselNext className="top-1/2 right-2 z-10 -translate-y-1/2" />
		</Carousel>
	) : (
		<FileDisplay username={username} post={{ postType, postOwner } as ScrapeResponse} file={files[0]} />
	);
}

export function FileDisplay({
	username,
	file,
	post: { postType, postOwner },
}: {
	username: string;
	file: string;
	post: ScrapeResponse;
}) {
	const url = `/api/storage/${username}/${postTypeString(postType)}/${postOwner}/${file}`;
	if (/\.(jpg)|(jpeg)|(webp)|(heic)$/.test(file)) {
		return <img src={url} />;
	} else if (/\.(mp4)|(webm)$/.test(file)) {
		return <video src={url} />;
	} else {
		return <a href={url}>{url}</a>;
	}
}
