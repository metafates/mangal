import createClient from "openapi-fetch";
import { type paths } from "./mangal";

const client = createClient<paths>({ baseUrl: '/api' })

export default client