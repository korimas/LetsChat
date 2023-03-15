export const config = {
    runtime: "edge",
};

const handler = async (req: Request): Promise<Response> => {
    let auth = false;
    if (process.env.PASSWORD) {
        auth = true;
    }
    return new Response(`{"authRequire": ${auth}}`);
};

export default handler;
