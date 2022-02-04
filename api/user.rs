use std::error::Error;

use http::{Response, StatusCode};
use vercel_lambda::{error::VercelError, lambda, IntoResponse, Request};

fn handler(_: Request) -> Result<impl IntoResponse, VercelError> {
    let response = Response::builder()
        .status(StatusCode::OK)
        .header("Content-Type", "text/plain")
        .body("Hello World")
        .expect("Internal Server Error");

    Ok(response)
}

// Start the runtime with the handler
fn main() -> Result<(), Box<dyn Error>> {
    Ok(lambda!(handler))
}
