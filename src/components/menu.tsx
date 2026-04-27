import { UserKeyIcon, DatabaseSearchIcon } from "lucide-react";

import {
	Sidebar,
	SidebarGroup,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarProvider,
} from "@/components/ui/sidebar";
import { InstagramIcon } from "@/components/ui/svgs/instagram";
import { TikTokIcon } from "@/components/ui/svgs/tiktok";
import { VSCOIcon } from "@/components/ui/svgs/vsco";

export function Menu() {
	return (
		<SidebarProvider className="flex flex-col">
			<Sidebar className="top-(--header-height) h-[calc(100svh-var(--header-height))]!">
				<SidebarGroup>
					<SidebarMenu>
						<SidebarMenuItem>
							<SidebarMenuButton>
								<UserKeyIcon className="h-4 w-4" />
								Authentication
							</SidebarMenuButton>
						</SidebarMenuItem>
						<SidebarMenuItem>
							<SidebarMenuButton>
								<InstagramIcon className="h-4 w-4" />
								Instagram Post
							</SidebarMenuButton>
						</SidebarMenuItem>
						<SidebarMenuItem>
							<SidebarMenuButton>
								<InstagramIcon className="h-4 w-4" />
								Instagram Highlight
							</SidebarMenuButton>
						</SidebarMenuItem>
						<SidebarMenuItem>
							<SidebarMenuButton>
								<InstagramIcon className="h-4 w-4" />
								Instagram Story
							</SidebarMenuButton>
						</SidebarMenuItem>
						<SidebarMenuItem>
							<SidebarMenuButton>
								<TikTokIcon className="h-40 w-40" />
								TikTok Post
							</SidebarMenuButton>
						</SidebarMenuItem>
						<SidebarMenuItem>
							<SidebarMenuButton>
								<VSCOIcon className="w-4" />
								VSCO Post
							</SidebarMenuButton>
						</SidebarMenuItem>
						<SidebarMenuItem>
							<SidebarMenuButton>
								<DatabaseSearchIcon className="w-4" />
								History
							</SidebarMenuButton>
						</SidebarMenuItem>
					</SidebarMenu>
				</SidebarGroup>
			</Sidebar>
		</SidebarProvider>
	);
}
