import { MatchRoute } from "@tanstack/react-router";
import { DatabaseSearchIcon, SidebarIcon, UserKeyIcon } from "lucide-react";

import { RakerLogo } from "@/components/logo";
import { ThemeToggle } from "@/components/theme-toggle";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { InstagramIcon } from "@/components/ui/svgs/instagram";
import { SnapchatIcon } from "@/components/ui/svgs/snapchat";
import { TikTokIcon } from "@/components/ui/svgs/tiktok";
import { VSCOIcon } from "@/components/ui/svgs/vsco";

export default function Header({ toggleMenu }: { toggleMenu: () => void }) {
	return (
		<header className="fixed inset-x-0 top-0 z-50 flex w-full items-center border-b bg-background/95 backdrop-blur-lg">
			<div className="grid h-(--header-height) w-full grid-cols-[1fr_auto_1fr] items-center px-4">
				<div className="justify-self-start">
					<Button className="h-8 w-8" variant="outline" size="icon" onClick={toggleMenu}>
						<SidebarIcon />
					</Button>
				</div>
				<div className="justify-self-center">
					<div className="flex flex-row items-center *:mx-1">
						<RakerLogo withVersion />
						<Separator orientation="vertical" />
						<MatchRoute to="/">
							<UserKeyIcon className="h-4" />
						</MatchRoute>
						<MatchRoute to="/instagram">
							<InstagramIcon className="h-4" />
							<Label>Post</Label>
						</MatchRoute>
						<MatchRoute to="/highlight">
							<InstagramIcon className="h-4" />
							<Label>Highlight</Label>
						</MatchRoute>
						<MatchRoute to="/story">
							<InstagramIcon className="h-4" />
							<Label>Story</Label>
						</MatchRoute>
						<MatchRoute to="/tiktok">
							<TikTokIcon className="h-4" />
							<Label>Post</Label>
						</MatchRoute>
						<MatchRoute to="/snapchat">
							<SnapchatIcon className="h-4" />
							<Label>Highlight</Label>
						</MatchRoute>
						<MatchRoute to="/vsco">
							<VSCOIcon className="h-4" />
							<Label>Post</Label>
						</MatchRoute>
						<MatchRoute to="/history">
							<DatabaseSearchIcon className="h-4" />
						</MatchRoute>
					</div>
				</div>
				<div className="justify-self-end">
					<ThemeToggle />
				</div>
			</div>
		</header>
	);
}
