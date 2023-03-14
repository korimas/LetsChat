import { IncomingMessage, ServerResponse } from 'http';

export const config = {
    runtime: 'edge',
};


interface RequestBody {
    // 在这里定义你需要的请求体字段
    [key: string]: any;
}

function parseRequestBody(req: IncomingMessage): Promise<RequestBody> {
    return new Promise((resolve, reject) => {
        let body = '';
        req.on('data', (chunk) => {
            body += chunk.toString();
        });
        req.on('end', () => {
            try {
                const json = JSON.parse(body);
                resolve(json);
            } catch (error) {
                reject(error);
            }
        });
        req.on('error', (error) => {
            reject(error);
        });
    });
}

export default async function (req: IncomingMessage, res: ServerResponse) {
    try {
        const body = await parseRequestBody(req);
        const { name = 'World' } = body;
        res.writeHead(200, { 'Content-Type': 'text/plain' });
        res.end(`Hello ${name}!`);
    } catch (error) {
        res.writeHead(400, { 'Content-Type': 'text/plain' });
        res.end("error");
    }
}
