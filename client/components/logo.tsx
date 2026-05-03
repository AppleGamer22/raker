import { Badge } from "@/components/ui/badge";

export function RakerLogo({ withVersion }: { withVersion?: boolean }) {
	return (
		<div className="flex flex-row items-center">
			<img alt="Raker Logo" src="/raker.svg" className="w-6" />
			{withVersion && (
				<Badge className="ml-2">
					{import.meta.env.DEV ? <>&alpha;</> : <>{import.meta.env.VITE_GIT_TAG}</>}
				</Badge>
			)}
		</div>
	);
}
