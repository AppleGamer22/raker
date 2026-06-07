import { createFileRoute, stripSearchParams } from "@tanstack/react-router";
import { z } from "zod";

import { scrapeTikTok } from "@/buf/raker/v1/raker-RakerServer_connectquery";
import {
	ExtractorFormShell,
	ExtractorSwitchField,
	ExtractorTextField,
	useExtractorForm,
} from "@/components/extractor-form";

const tikTokSearchDefaults = {
	owner: "",
	post: "",
	incognito: false,
};

const tikTokSearchSchema = z.object({
	owner: z.string().catch(tikTokSearchDefaults.owner),
	post: z.string().catch(tikTokSearchDefaults.post),
	incognito: z.boolean().catch(tikTokSearchDefaults.incognito),
});

export const Route = createFileRoute("/tiktok")({
	component: TikTok,
	search: {
		middlewares: [stripSearchParams(tikTokSearchDefaults)],
	},
	validateSearch: tikTokSearchSchema,
});

function TikTok() {
	const { owner, post, incognito } = Route.useSearch();
	const navigate = Route.useNavigate();
	const { form, result, setResult, isPending } = useExtractorForm({
		navigate,
		search: { owner, post, incognito },
		validators: {
			onChange: tikTokSearchSchema,
			onSubmit: tikTokSearchSchema.extend({
				owner: z.string().min(1, "post owner is required"),
				post: z.string().min(1, "post ID is required"),
			}),
		},
		mutation: scrapeTikTok,
		autoSubmitWhen: ({ owner, post }) => owner.length > 0 || post.length > 0,
		buildMutationArgs: ({ owner, post, incognito }) => ({ owner, post, incognito }),
		buildSearch: ({ owner, post, incognito }) => ({ owner, post, incognito }),
	});

	return (
		<ExtractorFormShell form={form} isPending={isPending} result={result} setResult={setResult}>
			<form.Field name="owner">
				{(field) => (
					<ExtractorTextField field={field} label="owner" placeholder="https://tiktok.com/@OWNER/video/id" />
				)}
			</form.Field>
			<form.Field name="post">
				{(field) => (
					<ExtractorTextField
						field={field}
						label="post ID"
						placeholder="https://tiktok.com/@owner/video/ID"
					/>
				)}
			</form.Field>
			<form.Field name="incognito">
				{(field) => <ExtractorSwitchField field={field} label="Incognito" />}
			</form.Field>
		</ExtractorFormShell>
	);
}
