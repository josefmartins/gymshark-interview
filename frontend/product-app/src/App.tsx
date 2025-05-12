import { useEffect, useState } from 'react';

interface Product {
  id: string;
  name: string;
  package_sizes: number[];
}

interface PackageCalcResult {
  size: number;
  units: number;
}

const API_URL = process.env.REACT_APP_API_URL + '/v1/products';

function App() {
  const [products, setProducts] = useState<Product[]>([]);
  const [name, setName] = useState('');
  const [sizes, setSizes] = useState('');
  const [newPackageSize, setNewPackageSize] = useState<number | ''>('');

  const [selectedProductId, setSelectedProductId] = useState('');
  const [calcUnits, setCalcUnits] = useState('');
  const [calcResult, setCalcResult] = useState<PackageCalcResult[]>([]);

  const fetchProducts = async () => {
    try {
      const res = await fetch(API_URL);
      const data = await res.json();
      setProducts(data.Data || []);
    } catch (err) {
      console.error('Error fetching products:', err);
    }
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  const createProduct = async (e: React.FormEvent) => {
    e.preventDefault();

    const package_sizes = sizes
      .split(',')
      .map(s => parseInt(s.trim(), 10))
      .filter(n => !isNaN(n));

    if (!name || package_sizes.length === 0) {
      alert('Enter valid name and sizes');
      return;
    }

    try {
      const res = await fetch(API_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, package_sizes })
      });
      if (!res.ok) throw new Error('Create failed');
      setName('');
      setSizes('');
      fetchProducts();
    } catch (err) {
      console.error('Error creating product:', err);
      alert('Failed to create product');
    }
  };

  const deleteProduct = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this product?')) return;

    try {
      const res = await fetch(`${API_URL}/${id}`, { method: 'DELETE' });
      if (!res.ok) throw new Error('Delete failed');
      fetchProducts();
    } catch (err) {
      console.error('Error deleting product:', err);
      alert('Failed to delete product');
    }
  };

  const addPackageSize = async (productId: string, size: number) => {
    try {
      const res = await fetch(`${API_URL}/${productId}/packageSizes/${size}`, {
        method: 'POST'
      });
      if (!res.ok) throw new Error('Add package size failed');
      setNewPackageSize('');
      fetchProducts();
    } catch (err) {
      console.error('Error adding package size:', err);
      alert('Failed to add package size');
    }
  };

  const removePackageSize = async (productId: string, size: number) => {
    try {
      const res = await fetch(`${API_URL}/${productId}/packageSizes/${size}`, {
        method: 'DELETE'
      });
      if (!res.ok) throw new Error('Remove package size failed');
      fetchProducts();
    } catch (err) {
      console.error('Error removing package size:', err);
      alert('Failed to remove package size');
    }
  };

  const calculatePackages = async () => {
    if (!selectedProductId || !calcUnits) {
      alert('Select product and enter units');
      return;
    }
    try {
      const res = await fetch(`${API_URL}/${selectedProductId}/calculate/${calcUnits}`, {
        method: 'POST'
      });
      if (!res.ok) throw new Error('Calc failed');
      const data = await res.json();
      setCalcResult(data.packages || []);
    } catch (err) {
      console.error('Calc error:', err);
      alert('Calculation failed');
    }
  };

  return (
    <div style={{ padding: '2rem', fontFamily: 'sans-serif' }}>
      <h1>Product Packages</h1>
    <hr style={{ margin: '2rem 0' }} />
      <h2>Calculate Packaging</h2>
      <div style={{ marginBottom: '1rem' }}>
        <select
          value={selectedProductId}
          onChange={e => setSelectedProductId(e.target.value)}
        >
          <option value="">Select product</option>
          {products.map(p => (
            <option key={p.id} value={p.id}>
              {p.name} ({p.id})
            </option>
          ))}
        </select>{' '}
        <input
          type="number"
          placeholder="Units"
          value={calcUnits}
          onChange={e => setCalcUnits(e.target.value)}
        />{' '}
        <button onClick={calculatePackages}>Calculate</button>
      </div>

      {calcResult.length > 0 && (
        <table border={1} cellPadding={8} style={{ borderCollapse: 'collapse', width: '30%' }}>
          <thead>
            <tr>
              <th>Package Size</th>
              <th>Units</th>
            </tr>
          </thead>
          <tbody>
            {calcResult.map((pkg, idx) => (
              <tr key={idx}>
                <td>{pkg.size}</td>
                <td>{pkg.units}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      <hr style={{ margin: '2rem 0' }} />

      <h2>Product List</h2>
      <table border={1} cellPadding={8} style={{ borderCollapse: 'collapse', width: '100%' }}>
        <thead>
          <tr>
            <th>Name</th>
            <th>ID</th>
            <th>Package Sizes</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {products.map(p => (
            <tr key={p.id}>
              <td>{p.name}</td>
              <td>{p.id}</td>
              <td>
                <table>
                  <tbody>
                    {p.package_sizes.map(size => (
                      <tr key={size}>
                        <td>{size}</td>
                        <td>
                          <button onClick={() => removePackageSize(p.id, size)}>Remove</button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
                <div style={{ marginBottom: '0.5rem' }}>
                  <input
                    type="number"
                    placeholder="size"
                    value={newPackageSize === '' ? '' : newPackageSize}
                    onChange={e => setNewPackageSize(Number(e.target.value))}
                  />
                  <button
                    onClick={() =>
                      newPackageSize !== '' && addPackageSize(p.id, newPackageSize)
                    }
                  >
                    Add Package Size
                  </button>
                </div>
              </td>
              <td>
                
                <div>
                  <button onClick={() => deleteProduct(p.id)}>Delete Product</button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      <hr style={{ margin: '2rem 0' }} />
      <h2>Create Product</h2>
      <div style={{ marginBottom: '1rem' }}>
        <form onSubmit={createProduct} style={{ marginBottom: '2rem' }}>
          <input
            type="text"
            placeholder="Product name"
            value={name}
            onChange={e => setName(e.target.value)}
          />{' '}
          <input
            type="text"
            placeholder="Package sizes (e.g. 250,500)"
            value={sizes}
            onChange={e => setSizes(e.target.value)}
          />{' '}
          <button type="submit">Create</button>
        </form>
      </div>

      
    </div>
  );
}

export default App;