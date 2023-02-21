import { get } from "svelte/store";
import { userStore } from "../stores/store.user";
import type { LayoutServerLoad } from './$types'

// get `locals.user` and pass it to the `page` store
export const load: LayoutServerLoad = async ({locals}) => {
    console.log(`layouts: ${JSON.stringify(locals)}`)
    return {user: locals.user};
}
