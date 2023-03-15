export const config = {
    runtime: "edge",
};

const handler = async (req: Request): Promise<Response> => {
    if (req.method == "GET") {
        let auth = false;
        if (process.env.PASSWORD) {
            auth = true;
        }
        return new Response(JSON.stringify({
            success:true,
            result:{authRequire: auth}
        }));
    }

    if (req.method == "POST") {
        const body = await req.json();

        if (body.password == process.env.PASSWORD) {
            return new Response(JSON.stringify({
                success:true
            }));
        }
        return new Response(JSON.stringify({
            success:false
        }));

    }

    throw new Error("Not support method");
};

export default handler;
