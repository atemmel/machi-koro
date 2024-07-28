export const get = <T>(url: string) => {
  return new Promise<T>((resolve, reject) => {
    fetch(url, {
      method: "GET",
      headers: {
        Accept: "application/json",
      },
    })
      .then(async (r) => {
        if (r.ok) {
          resolve(await r.json());
        } else {
          reject(r.status);
        }
      })
      .catch(() => {
        reject(0);
      });
  });
};

export const post = <T>(url: string, payload: any = undefined) => {
  return new Promise<T>((resolve, reject) => {
    fetch(url, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: payload,
    })
      .then(async (r) => {
        if (r.ok) {
          resolve(await r.json());
        } else {
          reject(r.status);
        }
      })
      .catch(() => {
        reject(0);
      });
  });
};
